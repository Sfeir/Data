package com.example;

import java.util.ArrayList;

import com.fasterxml.jackson.annotation.JsonProperty;

import org.apache.avro.reflect.Nullable;
import org.apache.beam.sdk.coders.AvroCoder;
import org.apache.beam.sdk.coders.DefaultCoder;
import org.codehaus.jackson.annotate.JsonIgnore;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;

@DefaultCoder(AvroCoder.class)
@JsonIgnoreProperties(ignoreUnknown = true)
//@JsonSerialize(include=JsonSerialize.Inclusion.NON_NULL)
public class Tweet {  
	/*
	@JsonIgnoreProperties(ignoreUnknown = true)
    public static class Coordinates {
        private String type;
		public String getType() {
			return type;
		}
		public void setType(String t) {
			type = t;
		}
    }
	*/
	
	@JsonIgnoreProperties(ignoreUnknown = true)
	public static class Entities {
		@Nullable
        private ArrayList<Hashtag> hashtags;

		public ArrayList<Hashtag> getHashtags() {
			return hashtags;
		}

		public void setHashtags(ArrayList<Hashtag> hashtags) {
			this.hashtags = hashtags;
		}
    }
	
	@JsonIgnoreProperties(ignoreUnknown = true)
	public static class Hashtag {
        private String text;

		public String getText() {
			return text;
		}

		public void setText(String text) {
			this.text = text;
		}
    }
    
	@JsonIgnoreProperties(ignoreUnknown = true)
    public static class User {
        private long id;
        private String name;
        @JsonProperty("screen_name")
        private String screen_name;
		public long getId() {
			return id;
		}
		public void setId(long a) {
			id = a;
		}
		public String getName() {
			return name;
		}
		public void setName(String a) {
			name = a;
		}
		public String getScreenName() {
			return screen_name;
		}
		public void setScreenName(String screenName) {
			screen_name = screenName;
		}
    }
    
    private Entities entities;
    private String created_at;
    
	public long getId() {
		return id;
	}
	public void setId(long id) {
		this.id = id;
	}
	public String getIn_reply_to_screen_name() {
		return in_reply_to_screen_name;
	}
	public void setIn_reply_to_screen_name(String in_reply_to_screen_name) {
		this.in_reply_to_screen_name = in_reply_to_screen_name;
	}
	public long getIn_reply_to_status_id() {
		return in_reply_to_status_id;
	}
	public void setIn_reply_to_status_id(long in_reply_to_status_id) {
		this.in_reply_to_status_id = in_reply_to_status_id;
	}
	public long getIn_reply_to_user_id() {
		return in_reply_to_user_id;
	}
	public void setIn_reply_to_user_id(long in_reply_to_user_id) {
		this.in_reply_to_user_id = in_reply_to_user_id;
	}
	public String getLang() {
		return lang;
	}
	public void setLang(String lang) {
		this.lang = lang;
	}
	public int getRetweet_count() {
		return retweet_count;
	}
	public void setRetweet_count(int retweet_count) {
		this.retweet_count = retweet_count;
	}
	public String getSource() {
		return source;
	}
	public void setSource(String source) {
		this.source = source;
	}
	public String getText() {
		return text;
	}
	public void setText(String text) {
		this.text = text;
	}
	public boolean isTruncated() {
		return truncated;
	}
	public void setTruncated(boolean truncated) {
		this.truncated = truncated;
	}
	public User getUser() {
		return user;
	}
	public void setUser(User user) {
		this.user = user;
	}
	
    private int favorite_count;
    private boolean favorited;
    
    private String filter_level;
    private long id;
    
    private String in_reply_to_screen_name;
    
    private long in_reply_to_status_id;
    
    private long in_reply_to_user_id;
    private String lang;
    
    private int retweet_count;
    private String source;
    private String text;
    private boolean truncated;
    private User user;
	public String getCreated_at() {
		return created_at;
	}
	public void setCreated_at(String created_at) {
		this.created_at = created_at;
	}
	
	public int getFavorite_count() {
		return favorite_count;
	}
	public void setFavorite_count(int favorite_count) {
		this.favorite_count = favorite_count;
	}
	public boolean isFavorited() {
		return favorited;
	}
	public void setFavorited(boolean favorited) {
		this.favorited = favorited;
	}
	public String getFilter_level() {
		return filter_level;
	}
	public void setFilter_level(String filter_level) {
		this.filter_level = filter_level;
	}
	public Entities getEntities() {
		return entities;
	}
	public void setEntities(Entities entities) {
		this.entities = entities;
	}
}
