#!/usr/bin/env bash
# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_ID="sfeir-data"
DATASET_ID="sfeir_opinions"
GCS_BUCKET="sfeir-opinion-analysis"

mvn compile exec:java \
  -Dexec.mainClass=com.sfeir.sentiment.IndexerPipeline\
  -Dexec.args="--project=$PROJECT_ID \
    --stagingLocation=gs://$GCS_BUCKET/staging/ \
    --runner=DataflowRunner \
    --tempLocation=gs://$GCS_BUCKET/temp/ \
    --maxNumWorkers=100 \
    --workerMachineType=n1-standard-2 \
    --zone=us-central1-f \
    --diskSizeGb=400 \
    --inputFile=gs://$GCS_BUCKET/input/*.txt \
    --bigQueryDataset=$DATASET_ID \
    --writeTruncate=false \
    --processedUrlHistorySec=130000 \
    --appName=Sentiments \
    --dedupeText=false \
    --saveProfilesToGcs=gs://$GCS_BUCKET/temp/ \
    --gcpTempLocation=gs://$GCS_BUCKET/temp/"

