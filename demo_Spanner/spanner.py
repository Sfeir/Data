from google.cloud import spanner
import matplotlib.pyplot as plt

x=[]
y=[]

def get_database(instance_id, database_id):
	spanner_client = spanner.Client()
	instance = spanner_client.instance(instance_id)
	return instance.database(database_id)

def query(database,query):
	with database.snapshot() as snapshot:
		return snapshot.execute_sql(query)

hashtag='GameOfThrones'
user='GOT Newz'

# top hashtags from the tweets
query1='select hashtag, count(*) from tweet_by_hashtag group by hashtag order by 2 desc limit 1' 
# number of tweets per day
query2='select substr(text,0,10) as date, count(*) as nb_tweets from tweet where text!="Vendredi" group by 1 order by substr(date,4,7)'
# nb of tweets containing the # @hashtag through time
query3 = 'select substr(text,0,10) as time, count(*) from tweet t, (select tweet_id from tweet_by_hashtag h where h.hashtag=\''+hashtag+'\') id where t.tweet_id=id.tweet_id  group by 1 order by substr(time,4,7)'
# top got tweeters
query4='select username, count(*) from tweet_by_user group by username order by 2 desc limit 20' 
# nb of tweets written by the user @user through time
query5 = 'select substr(text,0,10) as time, count(*) from tweet t, (select tweet_id from tweet_by_user u where u.username=\''+user+'\') id where t.tweet_id=id.tweet_id  group by 1 order by substr(time,4,7)'


#get_database("tweets","got")
results = query(get_database("tweets","got"),query5)
for row in results:
	#print row[0], row[1]
	x.append(row[0])
	y.append(row[1])

'''
x = x[-7:]+x[:-7]
y = y[-7:]+y[:-7]
'''

plt.figure()
plt.plot(range(len(x)),y)
plt.xticks(range(len(x)),x,rotation=90)
plt.xlabel("Day")
plt.ylabel("Number of tweets")
plt.title("Number of tweets per day from the user '"+user+"'")
#plt.title("Number of tweets per day containg hashtag #"+hashtag)
#plt.savefig("nbtweets_per_day")
plt.show()


