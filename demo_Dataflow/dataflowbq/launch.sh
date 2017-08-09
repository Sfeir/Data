#!/bin/bash


project=sfeir-data
tmpPath=gs://testdataflow123456789/pubsub_to_bq/dataflow/tmp
template=gs://testdataflow123456789/pubsub_to_bq/dataflow/template
dataset=SeriesDataflow
table=GOT


subscription=projects/${project}/subscriptions/subname_tuto2
tableId=${project}:${dataset}.${table}

mvn compile exec:java \
-Dexec.mainClass=com.example.PubSubToBigQueryJob \
-Dexec.args="\
--runner=DataflowRunner \
--streaming=true \
--jobName=tweets-pubsub-to-bq \
--project=$project \
--zone=europe-west1-b \
--gcpTempLocation=$tmpPath \
--workerMachineType=n1-standard-1 \
--diskSizeGb=20 \
--numWorkers=1 \
--maxNumWorkers=3 \
--autoscalingAlgorithm=THROUGHPUT_BASED \
--subscription=$subscription \
--table=$tableId" \



