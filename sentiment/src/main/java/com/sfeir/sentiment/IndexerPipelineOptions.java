/*******************************************************************************
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *******************************************************************************/
package com.sfeir.sentiment;

import org.apache.beam.runners.dataflow.options.DataflowPipelineOptions;
import org.apache.beam.sdk.options.Default;
import org.apache.beam.sdk.options.Description;

/**
 * Options supported by {@link IndexerPipeline}.
 */
public interface IndexerPipelineOptions extends DataflowPipelineOptions {

    @Description("Source files. Use with RecordFile and GDELTbucket sources. Names can be a pattern like *.gz or *.txt")
    String getInputFile();

    void setInputFile(String value);

    @Description("Sink BigQuery dataset name")
    String getBigQueryDataset();

    void setBigQueryDataset(String dataset);

    @Description("Deduplicate based on text content")
    Boolean getDedupeText();

    void setDedupeText(Boolean dedupeText);

    @Description("Text indexing option: Index as Short text or as an Article. Default: false")
    @Default.Boolean(false)
    Boolean getIndexAsShorttext();

    void setIndexAsShorttext(Boolean indexAsShorttext);

    @Description("Truncate sink BigQuery dataset before writing")
    @Default.Boolean(false)
    Boolean getWriteTruncate();

    void setWriteTruncate(Boolean writeTruncate);

    @Description("Time window in seconds counting from Now() for processed url cache. Set to " + Integer.MAX_VALUE + " to cache all Urls, don't set for no cache.")
    Integer getProcessedUrlHistorySec();

    void setProcessedUrlHistorySec(Integer seconds);

}
