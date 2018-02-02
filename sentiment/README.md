# Sample: Opinion Analysis of News
This sample uses Cloud Dataflow to build an opinion analysis processing pipeline for
news.

Opinion Analysis can be used for lead generation purposes, user research, or 
automated testimonial harvesting.

## About the sample

This sample contains:

* Cloud Dataflow pipelines for importing bounded (batch) raw data from files in Google Cloud Storage

## How to run the sample
The steps for configuring and running this sample are as follows:

- Setup your Google Cloud Platform project and permissions.
- Install tools necessary for compiling and deploying the code in this sample.
- Create or verify a configuration for your project.
- Clone the sample code
- Create the BigQuery dataset
- Deploy the Dataflow pipelines
- Clean up

### Prerequisites

Setup your Google Cloud Platform project and permissions

* Select or Create a Google Cloud Platform project.
  In the [Google Cloud Console](https://console.cloud.google.com/project), select
  **Create Project**.

* [Enable billing](https://support.google.com/cloud/answer/6293499#enable-billing) for your project.

* [Enable](https://console.cloud.google.com/flows/enableapi?apiid=dataflow,compute_component,logging,storage_component,storage_api,bigquery,pubsub,datastore) the Google Dataflow, Google Cloud Storage, and other APIs necessary to run the example. 

Install tools necessary for compiling and deploying the code in this sample, if not already on your system, specifically git, Google Cloud SDK, Java and Maven (for Dataflow pipelines):

* Install [`git`](https://git-scm.com/downloads).

* [Download and install the Google Cloud SDK](http://cloud.google.com/sdk/).

* Download and install the [Java Development Kit (JDK)](http://www.oracle.com/technetwork/java/javase/downloads/index.html) version 1.8 or later. Verify that the JAVA_HOME environment variable is set and points to your JDK installation.

* [Download](http://maven.apache.org/download.cgi) and [install](http://maven.apache.org/install.html) Apache Maven.


Create and setup a Cloud Storage bucket and Cloud Pub/Sub topics

* [Create a Cloud Storage bucket](https://console.cloud.google.com/storage/browser) for your project. This bucket will be used for staging your code, as well as for temporary input/output files. For consistency with this sample, select Multi-Regional storage class and United States location.

* Create folders in this bucket `staging`, `input`, `temp`

Create or verify a configuration for your project

* Authenticate with the Cloud Platform. Run the following command to get [Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials).

  `gcloud auth application-default login`

* Create a new configuration for your project if it does not exist already

  `gcloud init`

* Verify your configurations

  `gcloud config configurations list`


Important: This tutorial uses several billable components of Google Cloud
Platform. New Cloud Platform users may be eligible for a [free trial](http://cloud.google.com/free-trial).



### Clone the sample code

To clone the GitHub repository to your computer, run the following command:

```
git clone https://github.com/Sfeir/Data/tree/master/sentiment
```

Go to the `sentiment` directory. The exact path
depends on where you placed the directory when you cloned the sample files from
GitHub.

```
cd sentiment
```

### Create the BigQuery dataset

* Make sure you've activated the gcloud configuration for the project where you want to create your BigQuery dataset

  `gcloud config configurations activate <config-name>`

* In shell, go to the `bigquery` directory where the build scripts and schema files for BigQuery tables and views are located

  `cd bigquery`

* Run the `build_dataset.sh` script to create the dataset, tables, and views. The script will use the PROJECT_ID variable from your active gcloud configuration, and create a new dataset in BigQuery named 'sfeir_opinions'. In this dataset it will create several tables and views necessary for this sample.

  `./build_dataset.sh`

* [optional] Later on, if you make changes to the table schema or views, you can update the definitions of these objects by running update commands:

  `./build_tables.sh update`
  
Table schema definitions are located in the *Schema.json files in the `bigquery` directory.

### Deploy the Dataflow pipelines


#### Download and install the Sirocco sentiment analysis packages

Download and install [Sirocco](https://github.com/datancoffee/sirocco), a framework maintained by [@datancoffee](https://medium.com/@datancoffee).

* Download the latest [Sirocco Java framework](https://github.com/datancoffee/sirocco/releases/) 1.0.6 jar file.

* Download the latest [Sirocco model](https://github.com/datancoffee/sirocco-mo/releases/) 1.0.3 jar file.

* Go to the directory where the downloaded sirocco-sa-1.0.6.jar and sirocco-mo-1.0.3.jar files are located.

* Install the Sirocco framework in your local Maven repository

```
mvn install:install-file \
  -DgroupId=sirocco.sirocco-sa \
  -DartifactId=sirocco-sa \
  -Dpackaging=jar \
  -Dversion=1.0.6 \
  -Dfile=sirocco-sa-1.0.6.jar \
  -DgeneratePom=true
```

* Install the Sirocco model file in your local Maven repository

```
mvn install:install-file \
  -DgroupId=sirocco.sirocco-mo \
  -DartifactId=sirocco-mo \
  -Dpackaging=jar \
  -Dversion=1.0.3 \
  -Dfile=sirocco-mo-1.0.3.jar \
  -DgeneratePom=true
```

#### Build and Deploy your Indexer pipeline to Cloud Dataflow


* Go to the `sentiment/scripts` directory

* Edit the `run_indexerjob.sh` file in your favorite text editor, e.g. `nano`. Specifically, set the values of the variables used for parametarizing your control Dataflow pipeline. Set the values of the PROJECT_ID, DATASET_ID and other variables at the beginning of the shell script.

* Go back to the `sentiment` directory and run a command to deploy the control Dataflow pipeline to Cloud Dataflow.


```
cd ..
scripts/run_indexerjob.sh &
```

#### Run a verification job

You can use the included news articles (from Google's blogs) in the `src/test/resources/input` directory to run a test pipeline.

* Upload the files in the `src/test/resources/input` directory into the GCS `input` bucket. Use the [Cloud Storage browser](https://console.cloud.google.com/storage/browser) to find the `input` directory you created in Prerequisites. Then, upload all files from your local `src/test/resources/input` directory.

* Run the indexer pipeline from the [Dataflow Cloud Console](https://console.cloud.google.com/dataflow)

* Once the Dataflow job successfully finishes, you can review the data it will write into your target BigQuery dataset. Use the [BigQuery console](https://bigquery.cloud.google.com/) to review the dataset.

* Enter the following query to list new documents that were indexed by the Dataflow job. The sample query is using the Standard SQL dialect of BigQuery.

```
SELECT * FROM sfeir_opinions.sentiment 
ORDER BY DocumentTime DESC
LIMIT 100
```

### Clean up

Now that you have tested the sample, delete the cloud resources you created to
prevent further billing for them on your account.