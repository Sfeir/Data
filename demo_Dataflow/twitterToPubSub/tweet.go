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
/*
Source: 

{"contributors":null,"coordinates":{"coordinates":[0,0],"type":""},"created_at":"Thu Aug 03 14:53:14 +0000 2017","entities":{"hashtags":[{"indices":[114,120],"text":"GoTS7"}],"media":null,"url":null,"user_mentions":[{"id":1045360130,"id_str":"1045360130","indices":[3,15],"name":"chiquito","screen_name":"chiquitogif"}]},"favorite_count":0,"favorited":false,"filter_level":"low","id":893122592005128193,"in_reply_to_screen_name":"","in_reply_to_status_id":0,"in_reply_to_user_id":0,"lang":"es","place":null,"retweet_count":0,"source":"\u003ca href=\"http://twitter.com\" rel=\"nofollow\"\u003eTwitter Web Client\u003c/a\u003e","text":"RT @chiquitogif: SPOILER: \"Game of Thrones 4x07\"-- TYRION, el más grande, pese a medir medio centímetros cúbicos. #GoTS7 https://t.co/60CCc…","truncated":false,"user":{"id":76940667,"name":"Alberto P. Castaños","screen_name":"albertoperezc"}}

; line: 1, column: 199] (through reference chain: com.example.StationUpdate["entities"]->com.example.StationUpdate$Entities["url"])
*/






