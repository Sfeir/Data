## HOW TO -  BIGQUERY

The goal of this document is to show you how to use some tools of the Google Cloud SDK. More precisely, we will be using BigQuery, Compute Engine and Data Stduio. Moreover, we will use the GO language to build a program that retrieves tweets (via the Twitter API) and stores them in a BigQuery table.

# Google Cloud SDK

First of all, download the Google Cloud SDK available via this [link](https://cloud.google.com/sdk/docs/)

Then, as showed on this link :

$ ./google-cloud-sdk/install.sh

Open a new terminal so that the changes take effect and type :

$ gcloud init
$ gcloud components update

that allows ou to initialize your account and update the SDK.
Finally type : 

$ gcloud auth application-default login

that allows you to link your account for any future request.

At this stage, you have properly installed the Google Cloud SDK and linked your account and a project.

# GO 

Now to install the GO language type : 

$ sudo apt-get install golang

Then, to setup the paths necessary for GO ad to your .profile :

export GOROOT=/usr/local/go
export GOPATH=$HOME/go/src
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

Now, if you want to create a new GO program, just create a new folder in $HOME/go/src. For instance, hello. Then, create a new file in that folder, for example hello.go.
Type your code in that file (note that if you want is to be executable with a main, the package must be named main). Then, go to that folder and type :

$ go build

If there are no errors in your code, and there is a main function, this will create an executable named hello. You can launch it by typing by typing :

$ ./hello

We have not talked about the imports yet. In the GO language, you can import API or github projects. In order to use them, you have to type in the terminal :

$ go get cloud.google.com/go/bigquery 

in order to import the BigQuery API for the GO language . Note that to use it, you also have to write in your code :

import (‘’cloud.google.com/go/bigquery‘’)

Moreover, to import github projects, follow the same procedure. However, you must have git installed. To make sure it is the case type :

$ sudo apt-get install git

Now that, GO programs are clear you can launch the program retreiving tweets and pushing it to a BigQuery table. Copy the folder finalTwitterBigQuery to $HOME/go/src and follow the steps described before. Finally type :

$ ./finalTwitterBigQuery

When launching it, the program retrieves tweets according to the mask you precised in stream.go and pushes it to the BigQuery Table you indicated in stream.go. Note that you can create the BigQuery table in the code or manually via the [platform](https://console.cloud.google.com/) in the BigQuery section and by creating a DataSet and then a table.
Note that the credentials correspond to the Twitter API, you can get yours [here](https://apps.twitter.com/).


# COMPUTE ENGINE

The goal of this application is to retrieve tweets continuously. This is why Compute Engine is used, launching the GO program on Compute Engine will allow it to run as much as wanted.
First go [there](https://console.cloud.google.com/) and go to the Compute Engine Section and then in VM instances. There, create a new instance Once created, connect to it by clicking the SSH button. 
You are now on the Virtual Machine that is under Debian distribution. Hence, the setting up of the GO language works like described before. To upload files from the computer to the VM, you can click on the top right corner of the screen and click Upload file.

It is advised to work as root. For that type :

$ sudo -s

It is also important that you do not forget to type :

$ gcloud auth application-default login

Besides, do not forget to edit your .profile and make the go get for the imports.
Now let us suppose you put the finalTwitterBigQuery file in /root/go/src. To launch it you just have to type :

$ cd /root/go/src/finalTwitterBigQuery
$ ./finalTwitterBigQuery

But if you launch it on the VM and then exit it, your program will also stop. That is when crons are useful. They allow you to launch a command at a certain moment and/or at a given frequency. For that, we will use the crontab command : 

$ crontab -l			list of current crons
$ crontab -e			edit list of crons
$ crontab -r			erase cron list

For our usage, we just want to launch our GO program once. For that type :

$ crontab -e

and add that line :

47 8 27 7 4 cd /root/go/src/finalTwitterBigQuery && ./finalTwitterBigQuery > /root/go/log

which means that finalTwitterBigQuery will be launched on Thursday 27th of July at 08:47. Adapt the line to match the current date. Note that you can check the date with the command $ date. 

You can now exit the VM and wait the launch of your program, you can check if it is working by querying the BigQuery table you are filling and check if the number of rows is growing (select count(*) from dataset_name.table_name). 

When you want to stop the program, just open the VM (by clicking on the SSH button like before) and then type :

$ ps -aux

and then look for the line ./finalTwitterBigQuery and the matching pid number. Then type :

$ kill pid_number

# BIG QUERY and DATA STUDIO

Now that you have filled (or are filling) our BigQuery table. We will see how to use, query and showcase the stored data. BigQuery is a tool that allows querying data with SQL language requests.
To access BigQuery, click on the BigQuery section of the Google Cloud Platform. Then write your query on the top of the screen. Once written, type Run Query to see the results. For example if ou want to query the table named GOT from DataSet Series you can write :

SELECT count(*) from Series.GOT

that will give you the number of rows you have in this table.

There is a button « Show options » that allows you (among other things) to query without using cached results. That can be helpful if your table is still being filled.  On the same line as this button on the very right there is a button you can click on that will help you for your queries. Indeed, if there is an error in your query it well tell you what is wrong and if the query is correct it will tell you the amount of data it will process.

On the other hand, Data Studio allows to represent your data in an easy way. First go [there](https://www.google.com/analytics/data-studio/) and sign in for Data Studio. Once done, go to “data sources” which allows you to choose the data you want to represent. Click on BigQuery→My Projects →Project_name -> DataSet_name. At this stage, you can choose your table as a data source. Once chosen, click connect and it will open a window with every column you have on the table asking you how Data Studio should consider it, what it should do with it. For instance, for numeric fields you can aggregate them in a lot of different ways. Once finished, click on “create report” and you will arrive to the report. This is where you can plot a number of different things such as Time series, bar charts, pie charts… Once one of these selected you just have to choose the parameters you want to showcase and the job is done.

However, for more complicated graphs, you first have to query the table on BigQuery, save the view and finally connect the view in Data Studio and plot what you wanted to show.

For example, if you want to show the 10 most retweeted tweets in our table, first query in BigQuery:  

select Text, count(*) from Series.GOT group by Text order by 2 desc limit 10

then click on “Save View” and name it (for example 10_MAX_RT). Then, all you have to do is going to Data Studio and connect to the view you just created and (dor instance) plot a pie chart.

These are some other interesting queries on the table GOT :

1. select Text, count(*) from Series.GOT group by Text order by 2 desc

2. select count(*) as nb_tweets,Language from Series.GOT group by Language order by 1 desc

3. select count(*) as nb_tweets,Language from Series.GOT where Language != "en" group by Language order by 1 desc

4. select count(*) as nb_tweets,left(Date,10) from (select * from Series.GOT where upper(Text) contains("ARYA")) group by 2 

The first one allows to display the most popular tweets (the ones that have been the most retweeted).
The second query returns the number of tweets per language, which allows to plot a very nice pie chart in Data Studio that displays the repartition of languages.
The third query is very similar to the second one as we have just removed the English language (because it was too much represented and thus was giving a less readable pie chart).
The last query returns the number of tweets containing “Arya” per day. We can see the evolution of the number of tweets on a character throughout the days via a time series.

Finally, BigQuery allows similarity requests. For instance, on our table we can wonder what words are the mst associated with a character. This query :


SELECT <br/>
&nbsp;&nbsp;&nbsp;&nbsp;upper(REGEXP\_EXTRACT(Text,r'\s+(\w*)\s+')) AS word, <br/>
&nbsp;&nbsp;&nbsp;&nbsp;COUNT(\*) AS count <br/>
FROM <br/>
  &nbsp;&nbsp;&nbsp;&nbsp; (select Text,count(\*) from (select \* from Series.GOT where upper(Text) contains("SAM")) group by Text order by 2 desc) <br/>
WHERE <br/>
&nbsp;&nbsp;&nbsp;&nbsp;not upper(REGEXP_EXTRACT(Text,r'\s+(\w*)\s+')) contains("GAME") and not upper(REGEXP\_EXTRACT(Text,r'\s+(\w\*)\s+')) contains("OF") and  <br/>&nbsp;&nbsp;&nbsp;&nbsp;not upper(REGEXP\_EXTRACT(Text,r'\s+(\w\*)\s+')) contains ("THRONES") and not upper(REGEXP\_EXTRACT(Text,r'\s+(\w\*)\s+')) contains("SAM") <br/>
GROUP BY 1 <br/>
ORDER BY count DESC <br/>
LIMIT 100; <br/>

allows to give the words that are the most associated with Sam. Indeed, this query has in its top results “TARLY” which is the last name of Sam in Game of Thrones.
