# DEMO - BIGTABLE

The goal is to learn how to use Big Table from the Google Cloud Platform. More precisely, we will be using tweets stored in a BigQuery table and store them in Big Table tables in order to be able to follow, through time, the tweets per hashtag and the tweets per user.

## BIG TABLE

Big Table is a very scalable and fast tool for querying huge databases. However, it is fast (around a millisecond for a query) when you choose the row key for your table wisely. It is often said that you build a table in Big Table for a specific query as you will choose the appropriate row key for one query. You can find more information about the schema design (here)[https://cloud.google.com/bigtable/docs/schema-design] and (here)[https://cloud.google.com/bigtable/docs/schema-design-time-series].

First, let's create the Big Table instance. Go (here)[https://console.cloud.google.com/bigtable/instances] and click on 'Create instance', put 'tweets' as instance name (and ID), choose Development Instance type (cheaper than production as you only have one node), choose 'europe-west1-c' for the zone, finally HDD as storage type (again it is cheaper) and type 'create'. Unlike Big Query or Spanner, you cannot create tables from the GCP console nor read or write data. For thoses purposes, you can do it from the terminal with the 'cbt' command or with API calls.

In order to create the table, I used the 'cbt' command. You can find a detailed quickstart for the 'cbt' command (here)[https://cloud.google.com/bigtable/docs/quickstart-cbt]. But in a nutshell, just type :
```bash
$ gcloud components update
$ gcloud components install cbt
```

Then, in your home directory create/modify the ~/.cbtrc file and type :
```bash
project = sfeir-data
instance = tweets
```
Once done, create the table 'got' :
```bash
$ cbt createtable got
```
You can check if the creation was successfull :
```bash
$ cbt ls
```
You should see 'got' printed.
Finally, create a column family :
```bash
$ cbt createfamily got data
```

At this stage, we have created the Big Table instance with a table with one column family. Now, in order to transfer the tweets from Big Query to Big Table we will use a Dataflow job.
		
## DATAFLOW

The code is in dataflowbqbt and to launch it :
```bash
$ cd dataflowbqbt
$ bash launch.sh
```
The job launches a query on a Big Query table and stores the result in a Big Table table, moreover we build a key to be able to filter the users per date.
You can check if everything is going fine by checking the dataflow job (here)[https://console.cloud.google.com/dataflow].

Once the job succeeded you can check if the rows have been inserted by typing :
```bash
$ cbt count got
```
which displays the number of rows in the table 'got'.

Finally, you can display the tweets of a user at a certain date by querying the Big Table table via the 'cbt' command:
```bash
$ cbt read got prefix="Nananona"
```

## BIG QUERY

Another way of querying data from a Big Table table is from Big Query using an external table (also known as "Federated table"). To do so go in the Big Query section of the GCP console and click on "create new table". For Location choose "Google Cloud Bigtable" and put as source "https://googleapis.com/bigtable/projects/sfeir-data/instances/tweets/tables/got". Then choose a dataset destination and put a table name. In column families  click on "Edit as Text" and paste the content of the file "GOT\_bt\_table\_description". Finally, check the two boxes at the end. Finally, click on "create table". Then, you can query this table which will actually query Big Table. For example if you want the number of tweets per user (in standard SQL) :
```sql
SELECT
  SUBSTR(rowkey,0,LENGTH(rowkey)-15),
  COUNT(*)
FROM
  `SeriesDataflow.GOT_bt` a,
  UNNEST(a.data.username.cell) b
WHERE
  b.value IS NOT NULL
GROUP BY
  1
ORDER BY
  2 DESC
```

