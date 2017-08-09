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
	"time"
	"math/rand"
	//"reflect"

	"github.com/antonholmquist/jason"
	"github.com/nhjk/oauth"
	"golang.org/x/net/context"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
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
						t.Coordinates = Coordinates{}
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
	
	ctx := context.Background()
	/*
	// [START auth]
	proj := os.Getenv("sfeir-data")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "sfeir-data environment variable must be set.\n")
		os.Exit(1)
	}
	*/
	client, err := pubsub.NewClient(ctx, "sfeir-data")
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	// [END auth]

	// List all the topics from the project.
	fmt.Println("Listing all topics from the project:")
	topics, err := list(client)
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}
	for _, t := range topics {
		fmt.Println(t)
	}
	
	const topic = "tweets"
	// Create a new topic called example-topic.
	//if err := create(client, topic); err != nil {
	//	log.Fatalf("Failed to create a topic: %v", err)
	//}
	

	
	// Get a Topic called topic
	t := client.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
	    log.Fatalf("Topic error: %v", err)
	}
	if !ok {
	    log.Fatalf("Topic doesn t exist: %v", err)
	}
	
	//sub, err := client.CreateSubscription(ctx, "subname_tuto2", pubsub.SubscriptionConfig{Topic: t})

	//sub := client.Subscription("subname_tuto2")

	
 	if err != nil {
	// Handle error.
 	}
	/*
	// Delete the topic.
	if err := delete2(client, topic); err != nil {
		log.Fatalf("Failed to delete the topic: %v", err)
	}
	*/
	t.Stop()

	// Add your credentials for Twitter API
	SetCredentials("pH6j63LWGbk7YCf4ni1tP5gGX", "CnwIPrKP9FqJ6c9T97mLee5JSoE0jlaSXFsLw1fASXjcuLN3Fh", "887246877657313280-YXjfva2PErcdhu4bKB4SbOgRdhP0USv", "vrPP1j2JfPTljCh2YrZ0c7WknrNG4jvhZUvRljfRxix4r")

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
	//text := ""
	for tweet := range tweets {
 
		jsonBytes, err := json.Marshal(tweet)
		if err != nil {
			panic(err)
		}

		// the model used to parse the string in date
		layOut := "Mon Jan 02 15:04:05 -0700 2006"

		// parse from string to a date			
		dateStamp, err := time.Parse(layOut,tweet.CreatedAt)

		if err != nil {
	        	fmt.Println(err)
	 	}
		
		// probability to create data late
		if (rand.Int31n(100)<=5){
			// create the illusion it was created 4 minutes ago in order to create data late 		
			dateStamp = dateStamp.Add(-4 * time.Minute)
	
		} 
		

		unixTime := dateStamp.Unix() * 1000	

		// Create the message with real timestamp
		msg := pubsub.Message{
			Data: []byte(jsonBytes),
			Attributes: map[string]string{
				"timestamp": strconv.FormatInt(unixTime, 10),
			},
		}	
		
		//text += tweet.User.Name + " $$$**$$$ " + tweet.Text + tweet.CreatedAt + " $$$**$$$ " + strconv.Itoa(len(tweet.Entities.Hashtags))  + " $$$**$$$ " + tweet.Lang
		// Publish a text message on the created topic.
		
		if err := publish(client, topic, msg); err != nil {
			log.Fatalf("Failed to publish: %v", err)
		}
		
		/*
		err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
	 		log.Printf("Got message: %s %s", m.Data, m.Attributes["timestamp"])
	 		m.Ack()
	 	})
		*/

		//fmt.Printf(text)
		/*
		//s := strconv.FormatFloat(3.1415, 'E', -1, 64)
		x := strconv.FormatFloat(tweet.Coordinates.Coordinates[0], 'f', -1, 64)
		y := strconv.FormatFloat(tweet.Coordinates.Coordinates[1], 'f', -1, 64)
		z := strconv.Itoa(len(tweet.Contributors))
		zz := strconv.Itoa(len(tweet.Entities.Hashtags))
		log.Println(tweet.User.Name, "tweeted", tweet.Text, "coords", x, y,"nb contributors", z, "date", tweet.CreatedAt, "with Hashtags", zz, "number of fav", string(tweet.FavoriteCount), "filterLevel", tweet.FilterLevel, "in reply to", tweet.InReplyToScreenName, "with language : ", tweet.Lang, "nb of RT", tweet.RetweetCount, "source", tweet.Source)
		*/	
	}
}

// Comes from pubsub


func create(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START create_topic]
	t, err := client.CreateTopic(ctx, topic)
	if err != nil {
		return err
	}
	fmt.Printf("Topic created: %v\n", t)
	// [END create_topic]
	return nil
}

func list(client *pubsub.Client) ([]*pubsub.Topic, error) {
	ctx := context.Background()

	// [START list_topics]
	var topics []*pubsub.Topic

	it := client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
	// [END list_topics]
}

func listSubscriptions(client *pubsub.Client, topicID string) ([]*pubsub.Subscription, error) {
	ctx := context.Background()

	// [START list_topic_subscriptions]
	var subs []*pubsub.Subscription

	it := client.Topic(topicID).Subscriptions(ctx)
	for {
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	// [END list_topic_subscriptions]
	return subs, nil
}

func delete2(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START delete_topic]
	t := client.Topic(topic)
	if err := t.Delete(ctx); err != nil {
		return err
	}
	fmt.Printf("Deleted topic: %v\n", t)
	// [END delete_topic]
	return nil
}

func publish(client *pubsub.Client, topic string, msg pubsub.Message) error {
	ctx := context.Background()
	// [START publish]
	t := client.Topic(topic)
	result := t.Publish(ctx, &msg)
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	// [END publish]
	return nil
}

func publishWithSettings(client *pubsub.Client, topic string, msg []byte) error {
	ctx := context.Background()
	// [START publish_settings]
	t := client.Topic(topic)
	t.PublishSettings = pubsub.PublishSettings{
		ByteThreshold:  5000,
		CountThreshold: 10,
		DelayThreshold: 100 * time.Millisecond,
	}
	result := t.Publish(ctx, &pubsub.Message{Data: msg})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	// [END publish_settings]
	return nil
}

func publishSingleGoroutine(client *pubsub.Client, topic string, msg []byte) error {
	ctx := context.Background()
	// [START publish_single_goroutine]
	t := client.Topic(topic)
	t.PublishSettings = pubsub.PublishSettings{
		NumGoroutines: 1,
	}
	result := t.Publish(ctx, &pubsub.Message{Data: msg})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	// [END publish_single_goroutine]
	return nil
}

func getPolicy(c *pubsub.Client, topicName string) (*iam.Policy, error) {
	ctx := context.Background()

	// [START pubsub_get_topic_policy]
	policy, err := c.Topic(topicName).IAM().Policy(ctx)
	if err != nil {
		return nil, err
	}
	for _, role := range policy.Roles() {
		log.Print(policy.Members(role))
	}
	// [END pubsub_get_topic_policy]
	return policy, nil
}

func addUsers(c *pubsub.Client, topicName string) error {
	ctx := context.Background()

	// [START pubsub_set_topic_policy]
	topic := c.Topic(topicName)
	policy, err := topic.IAM().Policy(ctx)
	if err != nil {
		return err
	}
	// Other valid prefixes are "serviceAccount:", "user:"
	// See the documentation for more values.
	policy.Add(iam.AllUsers, iam.Viewer)
	policy.Add("group:cloud-logs@google.com", iam.Editor)
	if err := topic.IAM().SetPolicy(ctx, policy); err != nil {
		log.Fatalf("SetPolicy: %v", err)
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	// [END pubsub_set_topic_policy]
	return nil
}

func testPermissions(c *pubsub.Client, topicName string) ([]string, error) {
	ctx := context.Background()

	// [START pubsub_test_topic_permissions]
	topic := c.Topic(topicName)
	perms, err := topic.IAM().TestPermissions(ctx, []string{
		"pubsub.topics.publish",
		"pubsub.topics.update",
	})
	if err != nil {
		return nil, err
	}
	for _, perm := range perms {
		log.Printf("Allowed: %v", perm)
	}
	// [END pubsub_test_topic_permissions]
	return perms, nil
}
