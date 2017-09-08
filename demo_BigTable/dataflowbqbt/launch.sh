#!/bin/bash


project=sfeir-data
tmpPath=gs://testdataflow123456789/bq_to_bt/dataflow/tmp
template=gs://testdataflow123456789/bq_to_bt/dataflow/template
dataset=SeriesDataflow
table=GOT


subscription=projects/${project}/subscriptions/subname_tuto2
tableId=${project}:${dataset}.${table}

mvn compile exec:java \
-Dexec.mainClass=com.example.BigQueryToBigTableJob \
-Dexec.args="\
--runner=DataflowRunner \
--jobName=tweets-bq-to-bt \
--project=$project \
--zone=europe-west1-b \
--tempLocation=$tmpPath \
--workerMachineType=n1-standard-1 \
--diskSizeGb=20 \
--numWorkers=1 \
--maxNumWorkers=3 \
--autoscalingAlgorithm=THROUGHPUT_BASED \
--subscription=$subscription \
--table=$tableId" \



