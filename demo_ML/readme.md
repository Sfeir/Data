## HOW TO -  TENSOR FLOW - ML

The goal of this document is to show you how to use some tools of the Google Cloud SDK. More precisely, we will be using Dataprep (Beta), DataLab and the Google Cloud Natural Language API. The goal is to 

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

We will be using Storage in for Datalab. It is quite easy to use as you just have to create a bucket and then use it as if it were a local file on your computer (in terms of path).

# DATAPREP

Dataprep is a GCP service that allows you to clean data from BigQuery (or Storage) at a high level. Note that it is in closed Beta so you have to ask for Dataprep (it took less than a day for them to accept my request). Once you have access to Dataprep (it should now appear in the GCP console), it will ask you for a Storage Bucket. Then, you will be on the main Dataprep page where you can go to "Datasets" and select "Import Data". You can either import data from your computer, from Storage or from BigQuery. Once selected, click on "import and wrangle". This will take some time to load your data and show it. Now you can do several actions in order to "clean" your data. Note that above each column, you can visualize your data (grouped by occurence and counted). By clicking on the columns you see, you will be proposed some "recipes" depending on the type of data you have clicked on. For example, if it is a string, it will propose one recipe allowing you to keep rows that match some String. Just click on a recipe and edit it if you want or just add it. You can also add new columns from other columns. When you finished cleaning your data, click on "run job" and select the action you want to perform (write in a new BigQuery table, in storage, append data...) and click run. This will create a DataFlow job that you can see in the DataFlow section of the GCP console.

Note that I tried to perform a Dataprep cleaning on data located in a BigQuery table located in EU and that it failed giving the error : "Cannot read and write in different locations: source: EU, destination: US". 

# DATALAB

DataLab is an interactive tool designed to browse, analyze, transform and visualize data. It also allows you to create Machine Learning models on the GCP. Datalab is executed on Google Compute Engine and has the advantage of being able to connect easily to a lot of cloud services.

In order to set up DataLab you can follow the instructions [here](https://cloud.google.com/datalab/docs/quickstarts). I have worked on DataLab to retrieve some tweets from BigQuery, analyze them with the API Cloud Natural Language and match them (as good as possible) with an image from a public dataset from The Metropolitan Museum of Art. This dataset contains indormation about each image coming from a machine learning algorithm from Google.
