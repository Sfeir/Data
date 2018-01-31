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

import com.google.api.services.bigquery.model.TableRow;
import org.joda.time.DateTime;
import org.joda.time.DateTimeZone;
import org.joda.time.Instant;
import org.joda.time.format.DateTimeFormat;
import org.joda.time.format.DateTimeFormatter;
import sirocco.util.IdConverterUtils;

public class IndexerPipelineUtils {

    public static final String WEBRESOURCE_TABLE = "webresource";
    public static final String DOCUMENT_TABLE = "document";
    public static final String SENTIMENT_TABLE = "sentiment";

    // IDs of known document collections
    public static final String DOC_COL_ID_GDELT_BUCKET = "03";

    public static final DateTimeFormatter dateTimeFormatYMD_HMS = DateTimeFormat
            .forPattern("yyyy-MM-dd HH:mm:ss")
            .withZoneUTC();

    public static final DateTimeFormatter dateTimeFormatYMD_HMS_MSTZ = DateTimeFormat
            .forPattern("yyyy-MM-dd HH:mm:ss.SSS z");

    public static final DateTimeFormatter dateTimeFormatYMD = DateTimeFormat.forPattern("yyyy-MM-dd").withZoneUTC();

    public static final DateTimeFormatter dateTimeFormatYMD_T_HMS_Z = DateTimeFormat
            .forPattern("yyyyMMdd'T'HHmmss'Z'");

    // Integration with Cloud NLP
    public static final String CNLP_TAG_PREFIX = "cnlp::";

    // TODO: Move the date parsing functions to DateUtils
    public static DateTime parseDateString(DateTimeFormatter formatter, String s) {
        DateTime result = null;
        try {
            result = formatter.parseDateTime(s);
        } catch (Exception e) {
            result = null;
        }
        return result;
    }

    public static Long parseDateToLong(DateTimeFormatter formatter, String s) {
        Long result = null;
        DateTime date = parseDateString(formatter, s);
        if (date != null)
            result = date.getMillis();
        return result;
    }

    public static Long parseDateToLong(String s) {
        Long pt = parseDateToLong(dateTimeFormatYMD_HMS, s);
        if (pt == null)
            pt = parseDateToLong(dateTimeFormatYMD, s);
        return pt;
    }

    // same as in sirocco.utils.IdConverterUtils
    public static int getDateIdFromTimestamp(long millis) {
        int result;
        DateTime dt = new DateTime(millis, DateTimeZone.UTC);
        int year = dt.getYear();
        int month = dt.getMonthOfYear();
        int day = dt.getDayOfMonth();
        result = day + (100 * month) + (10000 * year);
        return result;
    }

    public static Integer parseDateToInteger(String s) {
        Integer result = null;
        Long l = parseDateToLong(s);
        if (l != null) {
            result = getDateIdFromTimestamp(l.longValue());
        }
        return result;
    }

    public static void setTableRowFieldIfNotNull(TableRow r, String field, Object value) {
        if (value != null)
            r.set(field, value);
    }

    public static String buildBigQueryProcessedUrlsQuery(IndexerPipelineOptions options) {
        String timeWindow = null;

        if (options.getProcessedUrlHistorySec() != null) {
            if (options.getProcessedUrlHistorySec() != Integer.MAX_VALUE) {
                Instant fromTime = Instant.now();
                fromTime = fromTime.minus(options.getProcessedUrlHistorySec() * 1000L);
                Integer fromDateId = IdConverterUtils.getDateIdFromTimestamp(fromTime.getMillis());
                timeWindow = "PublicationDateId >= " + fromDateId;
            }
        }

        if (timeWindow != null)
            timeWindow = "WHERE " + timeWindow;

        String result = "SELECT Url, MAX(ProcessingTime) AS ProcessingTime\n" + "FROM " + options.getBigQueryDataset()
                + "." + WEBRESOURCE_TABLE + "\n" + timeWindow + "\n" + "GROUP BY Url";

        return result;
    }

    public static String buildBigQueryProcessedDocsQuery(IndexerPipelineOptions options) {
        String timeWindow = null;

        if (options.getProcessedUrlHistorySec() != null) {
            if (options.getProcessedUrlHistorySec() != Integer.MAX_VALUE) {
                Instant fromTime = Instant.now();
                fromTime = fromTime.minus(options.getProcessedUrlHistorySec() * 1000L);
                Integer fromDateId = IdConverterUtils.getDateIdFromTimestamp(fromTime.getMillis());
                timeWindow = "PublicationDateId >= " + fromDateId;
            }
        }

        if (timeWindow != null)
            timeWindow = "WHERE " + timeWindow;

        String result = "SELECT DocumentHash, MAX(ProcessingTime) AS ProcessingTime\n" + "FROM " + options.getBigQueryDataset()
                + "." + DOCUMENT_TABLE + "\n" + timeWindow + "\n" + "GROUP BY DocumentHash";

        return result;
    }

    public static void validateIndexerPipelineOptions(IndexerPipelineOptions options) throws Exception {

        if (options.getInputFile().isEmpty())
            throw new IllegalArgumentException(
                    "Input file path pattern needs to be specified when GCS is being used as a source type.");

        if (options.getBigQueryDataset().isEmpty()) {
            throw new IllegalArgumentException("Sink BigQuery dataset needs to be specified.");
        }

        options.setStreaming(false);
    }

}
