## HOW TO -  DATAFLOW

The goal of this document is to show you how to use some tools of the Google Cloud SDK. More precisely, we will be using DataFlow, PubSub and Storage. Moreover, we will use the GO language to build a program that retrieves tweets (via the Twitter API) and sends them to a PubSub topic. Then we retrieve the sent data and will apply some specific functions to this data in streaming mode.

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

# STORAGE

We will be using Storage in for the Dataflow demonstrator. It is quite easy to use as you just have to create a bucket and then use it as if it were a local file on your computer (in terms of path).

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

Now that, GO programs are clear you can launch the program retreiving tweets and pushing it to a PubSub topic. Copy the folder twitterToPubSub to $HOME/go/src and follow the steps described before. Finally type :

$ ./twitterToPubSub

When launching it, the program retrieves tweets according to the mask you precised in stream.go and pushes it to the PubSub topic you indicated in stream.go (you can also create the topic in commented lines, and can also create a subscription). Note that you can create the PubSub topic and the subscription in the code or manually via the [platform](https://console.cloud.google.com/) in the PubSub section.
Note that the credentials correspond to the Twitter API, you can get yours [here](https://apps.twitter.com/).


# DATAFLOW

Here, we will be using the Java:SDK 2.X because (as I write these lines) the Python SDK does not handle streaming and therefore does not handle Data Late either (and those two parameters are of paramount importance in our demo).

First of all, there is a very good way to start using DataFlow with this SDK by following the steps [here](https://cloud.google.com/dataflow/docs/quickstarts/quickstart-java-maven). If you need help installing or setting up java environment on your computer there is a good guide [here](https://www.digitalocean.com/community/tutorials/how-to-install-java-with-apt-get-on-ubuntu-16-04).

After those steps, you have successfully started you first DataFlow job and checked if everything went as planned. Note that this example is in batch mode whereas we will be showcasing a streaming case.

The code corresponding to the DataFlow job in the file "dataflowbq". Note that for the job to run wihout any issues you should have created a PubSub topic, have subscribed to it and also have a table in BigQuery matching the "bigQueryTable_SeriesDataflow.GOT" description.

To launch the DataFlow job just go to dataflowbq/ and type "$ bash launch.sh". Note that you can edit the launch.sh file to  change the name of the BigQuery Dataset, table, the name of the project, the PubSub Bucket name.... These parameters will be available through the Options class defined at dataflowbq/src/main/java/com/example/Options.java.

Once launched, the pipeline will take approximately one minute to start working fine. You can either check it on the terminal or on the google cloud console online in the Dataflow section (just click on the name of your job and you will be able to see the acyclic graph corresponding to your pipeline). Each node of this graph represents one (or several) "apply" which correspond to a function you apply to the data you have collected or are collecting.

For example in our case you can see in "dataflowbq/src/main/java/com/example/PubSubToBigQueryJob.java" file in the main that we first retrieve the data from PubSub with :

PCollection<String> lines = p
				.apply("Read from PubSub",PubsubIO.readStrings()
                     		.withTimestampAttribute(TIMESTAMP_ATTRIBUTE)
                        	.fromSubscription(options.getSubscription()));

with p being the pipeline created just before. This ".apply" represents the first node of the acyclic graph.

After retrieving this data, we map the String (sent via the GO program) to a Tweet (corresponds to "dataflowbq/src/main/java/com/example/Tweet.java" class). Now, on the one hand we push this data to a BigQuery table and on the other hand we apply windowing to the same data. This step corresponds to the separation in two nodes. Windowing allows to make calculations (like wordcounts) on unbounded sets of data. Indeed, for such sets of data, you do not know when the data is "finished" (it might actually never) o you have to window it in (for instance) windows of 4 minutes. It means that your unbounded data will be cut in 4 minutes batches. The following code allows windowing : 

PCollection<Tweet> windowed = t1.apply("Windowing",Window.<Tweet>into(FixedWindows.of(Duration.standardMinutes(4))).withAllowedLateness(Duration.standardMinutes(5)).accumulatingFiredPanes());

As we can see, this returns fixed windowed data of 4 minutes each with an allowed lateness time of 5 minutes. This allowed lateness corresponds to the time after which data can still be considrered as being part of a window even if it is late : this is is what is called Data Late. For example, if a window opens at 9:00 then it will close at 9:04. However, if a data arrives at 9:06 but with a timestamp between 9:00 and 9:04 then it will be considered as Data Late from this window. However if a data arrives at 9:15 but with a timestamp between 9:00 and 9:04 it will be discarded as it is beyond the allowed lateness (after the closure of a window).

Moreover, Data Late can be dealt with by two main ways. The first one is by using ".accumulatingFiredPanes()" which means that when A Data LAte arrives it will be considered as part of the group with all the other data from its window (Late or not) whereas if you use ".discardingFiredPanes()" it will be considered as a single Data and will not be "merged" with the other data from its window.

Finally, after windowing we retrieve the Hashtags from the tweets and count them. Then, if it is not Data Late we write it to GCP's Storage and if it is Data Late we log it. You can see the logs by clicking on logs and then on Stackdriver.
                
