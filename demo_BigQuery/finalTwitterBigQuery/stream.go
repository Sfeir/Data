package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"strconv"

	"github.com/antonholmquist/jason"
	"github.com/nhjk/oauth"

        // Imports the Google Cloud BigQuery client package.
        "cloud.google.com/go/bigquery"
        "golang.org/x/net/context"
)

const (
	base   = "https://stream.twitter.com/1.1/statuses/"
	sample = "sample.json"
	filter = "filter.json"
)

var (
	// oauth credentials
	ClientKey    string
	ClientSecret string
	TokenKey     string
	TokenSecret  string

	// streaming api parameters
	FilterLevel string
	Language    []string // BCP 47 language identifiers.
	Track       []string // Tracking words and phrases.
	Follow      []string // User ids to follow.
	Locations   []string // longitude, latitude coordinates.
)

// SetCredentials is a helper function to set all oauth credentials at once
func SetCredentials(ck, cs, tk, ts string) {
	ClientKey, ClientSecret, TokenKey, TokenSecret = ck, cs, tk, ts
}

var s chan *Tweet

// Stream takes parameters to the twitter streaming api and returns a channel
// tweets matching those parameters
func Stream() (<-chan *Tweet, error) {
	req := endpoint()
	authorize(req)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("tweetstream: Non 200 status code: %s", resp.Status)
	}

	s := stream(resp)
	return s, nil
}

// Closes the stream
func Close() {
	close(s)
}

func authorize(req *http.Request) {
	c := &oauth.Consumer{ClientKey, ClientSecret}
	t := &oauth.Token{TokenKey, TokenSecret}
	c.Authorize(req, t)
}

func endpoint() (req *http.Request) {
	if Track == nil && Follow == nil && Locations == nil {
		req, _ = http.NewRequest("GET", base+sample, nil)
	} else {
		req, _ = http.NewRequest("POST", base+filter, nil)
	}

	values := req.URL.Query()

	if FilterLevel != "" {
		values.Set("filter_level", FilterLevel)
	}
	setIfNotEmpty(&values, "language", Language)
	setIfNotEmpty(&values, "track", Track)
	setIfNotEmpty(&values, "follow", Follow)
	setIfNotEmpty(&values, "locations", Locations)

	req.URL.RawQuery = values.Encode()

	return
}

func setIfNotEmpty(values *url.Values, name string, param []string) {
	if len(param) > 0 {
		values.Set(name, strings.Join(param, ","))
	}
}

func stream(resp *http.Response) <-chan *Tweet {
	tweetc := make(chan *Tweet)

	// filters tweets and locations in compliance with twitter's "delete" and
	// "scrub_geo" json messages
	delc, scrubc := make(chan int64), make(chan scrub)
	filteredTweetsc := filterTweets(tweetc, delc, scrubc)

	jsonc := toJsonc(resp.Body)

	go func() {
		defer close(tweetc)

		for j := range jsonc {
			rs := roots(j)

			switch {
			// a tweet (text is a guaranteed root for a tweet)
			case rs["text"]:
				tweet := new(Tweet)
				err := json.Unmarshal(j, tweet)
				if err != nil {
					log.Println("tweetsream: invalid json", j)
				} else {
					tweetc <- tweet
				}

			// message to remove a tweet with the given id
			case rs["delete"]:
				jas, err := jason.NewObjectFromBytes(j)
				if err != nil {
					log.Println("tweetstream: invalid json", j)
				} else {
					id, _ := jas.GetNumber("delete", "status", "id")
					a, err := strconv.ParseInt(string(id), 10, 64)
					if err != nil {
					    panic(err)
					}
					delc <- a
				}

			// message to scrub geolocation data from a range of tweets for a specified user
			case rs["scrub_geo"]:
				jas, err := jason.NewObjectFromBytes(j)
				if err != nil {
					log.Println("tweetstream: invalid json", j)
				} else {
					uid, _ := jas.GetNumber("scrub_geo", "user_id")
					tid, _ := jas.GetNumber("scrub_geo", "up_to_status_id")
					b, err := strconv.ParseInt(string(uid), 10, 64)
					if err != nil {
					    panic(err)
					}
					c, err := strconv.ParseInt(string(tid), 10, 64)
					if err != nil {
					    panic(err)
					}					
					scrubc <- scrub{b, c}
				}

			// indicates that a filtered stream has matched more more tweets than its
			// current rate limit allows to be delivered
			case rs["limit"]:
				jas, err := jason.NewObjectFromBytes(j)
				if err != nil {
					log.Println("tweetstream: invalid json", j)
				} else {
					unsentCount, _ := jas.GetNumber("limit", "track")
					d, err := strconv.ParseInt(string(unsentCount), 10, 64)
					if err != nil {
					    panic(err)
					}
					log.Println("tweetstream: rate limit reached,", d, "unsent tweets")
				}

			case rs["disconnect"]:
				// indicates that a status/user is witheld in a list of countries
				//case roots["status_withheld"], roots["user_withheld"], "":
			}
		}

	}()

	return filteredTweetsc
}

// filterTweets removes tweets and geolocation data from tweets in response to the
// "delete" and "scrub_geo" json messages sent by the twitter streaming api
func filterTweets(tweetc <-chan *Tweet, delc <-chan int64, scrubc <-chan scrub) <-chan *Tweet {
	delSet, scrubSet := make(map[int64]bool), make(map[scrub]bool)
	filteredTweetsc := make(chan *Tweet)

	go func() {

		for {
			select {
			case d := <-delc:
				delSet[d] = true

			case s := <-scrubc:
				scrubSet[s] = true

			case t := <-tweetc:
				// check if this tweet needs to be deleted
				for id, _ := range delSet {
					// remove delete notice if its expired
					if t.Id > id {
						delete(delSet, id)

						// don't send tweet
					} else if t.Id == id {
						break
					}
				}

				// check if this tweet needs to be scrubbed of geo data
				for s, _ := range scrubSet {
					// remove scrub notice if its expired
					if t.Id > s.upToTweetId {
						delete(scrubSet, s)

						// scrub geo data
					} else if t.User.Id == s.userId {
						//t.Coordinates = Coordinates{}
					}
				}

				// send tweet after passing through the 2 filters
				filteredTweetsc <- t
			}
		}
	}()

	return filteredTweetsc
}

type scrub struct {
	userId      int64
	upToTweetId int64
}

// roots returns the roots of a json object
func roots(p []byte) map[string]bool {
	rs := make(map[string]bool)

	var j map[string]interface{}
	json.Unmarshal(p, &j)

	for r, _ := range j {
		rs[r] = true
	}

	return rs
}

// toJsonc reads the body splitting on newlines into a chan of json []byte's
func toJsonc(body io.ReadCloser) <-chan []byte {
	jc := make(chan []byte)

	scanner := bufio.NewScanner(body)

	go func() {
		defer body.Close()

		var j []byte
		for scanner.Scan() {
			j = scanner.Bytes()
			jc <- j
		}

		close(jc)
	}()

	return jc
}

func main() {

	// Add your credentials
	SetCredentials("pH6j63LWGbk7YCf4ni1tP5gGX", "CnwIPrKP9FqJ6c9T97mLee5JSoE0jlaSXFsLw1fASXjcuLN3Fh", "887246877657313280-YXjfva2PErcdhu4bKB4SbOgRdhP0USv", "vrPP1j2JfPTljCh2YrZ0c7WknrNG4jvhZUvRljfRxix4r")

	ctx := context.Background()

        // Sets your Google Cloud Platform project ID.
        projectID := "sfeir-data"

        // Creates a client.
        client, err := bigquery.NewClient(ctx, projectID)
        if err != nil {
                log.Fatalf("Failed to create client: %v", err)
        }
	
	// reference to Dataset
	myDataset := client.Dataset("Series")


	// create dataset (if already created, will raise error)
	/*if err := myDataset.Create(ctx); err != nil {
    		log.Fatalf("Failed to create dataset: %v", err)
	}
	*/

	// reference to a table from the dataset above
	table := myDataset.Table("GOT")

	u := table.Uploader()
	// Item implements the ValueSaver interface.
	type Item2 struct {
		Username  string
		Text  string
		X float64
		Y float64
		Date string
		Nb_Hashtags int
		Filter_Level string
		Language string
		Nb_Retweets int
		Nb_Favorites int
	}


	// Add optional parameters. If none are set, it defaults to a sample stream.
	Track = []string{"Game of Thrones"}

	// Stream.
	tweets, _ := Stream()
	/*
	Contributors        []int64                `json:"contributors"`
	Coordinates         Coordinates            `json:"coordinates"`
	CreatedAt           string                 `json:"created_at"`
	Entities            Entities               `json:"entities"`
	FavoriteCount       int                    `json:"favorite_count"`
	Favorited           bool                   `json:"favorited"`
	FilterLevel         string                 `json:"filter_level"`
	Id                  int64                  `json:"id"`
	InReplyToScreenName string                 `json:"in_reply_to_screen_name"`
	InReplyToStatusId   int64                  `json:"in_reply_to_status_id"`
	InReplyToUserId     int64                  `json:"in_reply_to_user_id"`
	Lang                string                 `json:"lang"`
	Place               map[string]interface{} `json:"place"`
	RetweetCount        int                    `json:"retweet_count"`
	Source              string                 `json:"source"`
	Text                string                 `json:"text"`
	Truncated           bool                   `json:"truncated"`
	User                User                   `json:"user"`
	*/
	for tweet := range tweets {
		//s := strconv.FormatFloat(3.1415, 'E', -1, 64)
		x := tweet.Coordinates.Coordinates[0]
		y := tweet.Coordinates.Coordinates[1]
		//z := strconv.Itoa(len(tweet.Contributors))
		zz := len(tweet.Entities.Hashtags)
		//fmt.Println(tweet.User.Name, "tweeted", tweet.Text, "coords", x, y,"nb contributors", z, "date", tweet.CreatedAt, "with Hashtags", zz, "number of fav", string(tweet.FavoriteCount), "filterLevel", tweet.FilterLevel, "in reply to", tweet.InReplyToScreenName, "with language : ", tweet.Lang, "nb of RT", tweet.RetweetCount, "source", tweet.Source)

		
		items := []*Item2{
			{Username: tweet.User.Name, Text: tweet.Text, X: x, Y: y, Date: tweet.CreatedAt, Nb_Hashtags: zz, Filter_Level: tweet.FilterLevel, Language: tweet.Lang, Nb_Retweets: tweet.RetweetCount, Nb_Favorites: tweet.FavoriteCount},
			//{Username: "Peter Dinklage", Text: "Winter is coming", X: float64(0.0), Y: float64(0.0), Date: "Vendredi", Nb_Hashtags: 3, Filter_Level: "low", Language: "eng", Nb_Retweets: 0, Nb_Favorites: 0},
			//{Username: "Emilia Clarke", Text: "Dragons are awesome"},
		}
		if err := u.Put(ctx, items); err != nil {
			log.Fatalf("Failed to put : %v", err,zz)
		}
	
		}
}
