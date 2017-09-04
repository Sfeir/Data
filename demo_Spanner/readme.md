## DEMO - SPANNER

The goal is to learn how to use Spanner from the Google Cloud Platform. More precisely, we will be using tweets stored in a BigQuery table and store them in Spanner tables in order to be able to follow, through time, the tweets per hashtag and the tweets per user.

# BIGQUERY

The first step is to create several spanner tables from one BigQuery table. For that purpose, we will create the tables in BigQuery (via SQL requests) and then export them as .csv files and then import them in Spanner. We will create four tables from the original table. The declaratio of these tables ca nbe found at /tables_bq. In order to fill these tables I used the following queries (in Standard SQL) : 

Table tweet_spanner_main which is just an intermediate table that allows to extract the hashtags from the tweets and generate a unique id for each tweet :
```sql
SELECT 
Username, Date, Text, hashtags, row_number() over() id
FROM 
(
  SELECT
    regexp_extract_all(Text,
      r"#([a-zA-ZÀ-ÖØ-öø-ÿ0-9]+)") as hashtags,Text,Date,Username 
  FROM
    `Series.GOT`
  WHERE
    _PARTITIONTIME BETWEEN TIMESTAMP("2017-07-15") AND TIMESTAMP("2017-08-28")
)
```
To fill the table tweet :
```sql
select tweet_id, date, tweet
from `Series.tweet_spanner_main` 
```
To fill the table tweet_by_user :
```sql
select user, tweet_id
from `Series.tweet_spanner_main` 
```
To fill the table tweet_by_hashtag :
```sql
select hashtag, tweet_id
from `Series.tweet_spanner_main`, unnest (hashtags) as hashtag
```
and then, in order to remove duplicates (we cannot have any as spanner would not accept them) :
```sql
SELECT
  hashtag,tweet_id
FROM (
  SELECT *,
    ROW_NUMBER() OVER(PARTITION BY hashtag,tweet_id) AS dup
  FROM
    [Series.tweet_by_hashtag])
WHERE
  dup = 1
```
# SPANNER

Now that we have created the tables we wanted and stored them as .csv files we can import them to spanner. To do so, I followed the instructions (here)[https://www.youtube.com/watch?v=d-D_KgQ3L34].
First, create an instance and a dataset in Spanner. Once done, create three tables which DDL can be found at /ddl. Once created you can import the tables from BigQuery to spanner. First, download the three last tables created in BigQuery as .csv on your computer in the folder /spannerimport. From /spannerimport just type :
```bash
$ python import.py --instance_id=tweets --database_id=got --table_id=tweet_by_hashtag --batchsize=6000 --data_file=tweet_by_hashtag.csv --format_file=hashtag_table.fmt
```
For example, this will fill the tweet_by_hashtag table in spanner. The format_file describes the format of your table. Note that you have to remove the first line of your .csv files (name of the columns). The batchsize is the number of rows sent per batch. The maximum size is 100 MB or 20 000 mutations (the nuber of mutation being the number of rows X the number of columns). So if your table contains three columns you can put a batchsize of 6 666 (if it does not exceed 100 MB).

Now that the spanner tables are filled, we can query them and display the number of tweets for a hashtag or for a user throughout time. Five queries are written in the Python script spanner.py, queries 3 and 5 display what we want. To launch the script and display the graph, type :
```bash
$ python spanner.py
```
