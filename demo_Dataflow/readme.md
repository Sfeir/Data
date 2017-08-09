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

Here, we will be using the Java:SDK 2.X as the Python SDK does not handle streaming //TODO.
