package main

type Tweet struct {
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
}

type Coordinates struct {
	Coordinates [2]float64 `json:"coordinates"`
	Type        string     `json:"type"`
}

type Entities struct {
	Hashtags     []Hashtag                `json:"hashtags"`
	Media        []map[string]interface{} `json:"media"`
	Urls         []map[string]interface{} `json:"url"`
	UserMentions []map[string]interface{} `json:"user_mentions"`
}

type Hashtag struct {
	Indices []int  `json:"indices"`
	Text    string `json:"text"`
}

type User struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}
