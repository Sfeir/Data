package com.example;

import static org.apache.beam.sdk.io.gcp.bigquery.BigQueryIO.Write.CreateDisposition.CREATE_NEVER;
import static org.apache.beam.sdk.io.gcp.bigquery.BigQueryIO.Write.WriteDisposition.WRITE_APPEND;

import java.io.IOException;

import org.apache.beam.sdk.Pipeline;
import org.apache.beam.sdk.io.gcp.bigquery.BigQueryIO;
import org.apache.beam.sdk.io.gcp.bigquery.TableDestination;
import org.apache.beam.sdk.io.gcp.pubsub.PubsubIO;
import org.apache.beam.sdk.options.PipelineOptionsFactory;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.MapElements;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.transforms.SerializableFunction;
import org.apache.beam.sdk.transforms.windowing.FixedWindows;
import org.apache.beam.sdk.transforms.windowing.PaneInfo.Timing;
import org.apache.beam.sdk.transforms.windowing.Window;
import org.apache.beam.sdk.values.KV;
import org.apache.beam.sdk.values.PCollection;
import org.apache.beam.sdk.values.ValueInSingleWindow;
import org.joda.time.DateTimeZone;
import org.joda.time.Duration;
import org.joda.time.format.DateTimeFormatter;
import org.joda.time.format.ISODateTimeFormat;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.common.WriteOneFilePerWindow;
import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.api.services.bigquery.model.TableRow;

public class PubSubToBigQueryJob {

    private static final ObjectMapper MAPPER = new ObjectMapper();

    private static final String TIMESTAMP_ATTRIBUTE = "timestamp";
    
    private static final Logger logger = LoggerFactory.getLogger(PubSubToBigQueryJob.class);

    public static void main(String[] args) {
        // Get the job's options from the command line arguments.
        final Options options = PipelineOptionsFactory.fromArgs(args).withValidation().as(Options.class);

        // Create the pipeline.
        final Pipeline p = Pipeline.create(options);
        
        PCollection<String> lines = p
                //
                // Read station updates from Pub/Sub using custom timestamps.
                //
                .apply("Read from PubSub",PubsubIO.readStrings()
                        .withTimestampAttribute(TIMESTAMP_ATTRIBUTE)
                        .fromSubscription(options.getSubscription()));
                
        //PCollection<String> fixed_windowed_lines = lines.apply(
          //      Window.<String>into(FixedWindows.of(Duration.standardMinutes(1))));
                //
                // Map lines to StationUpdate.
                //
                //.apply("DecodeJson", MapElements.into(TypeDescriptor.of(Tweet.class)).via(PubSubToBigQueryJob::mapJsonToTweet))
                //PCollection<Integer> wordLengths = words
        //PCollection<Integer> l = lines.apply(MapElements.into(TypeDescriptors.integers()).via((String word) -> word.length()));
        //PCollection<Integer> l = lines.apply("func1",ParDo.of(new ComputeWordLengthFn()));
        //PCollection<Integer> l2 = lines
        		PCollection<Tweet> t1 = lines.apply("Map to Tweet",ParDo.of(new JsonMap()));
        		//.apply(TextIO.write().to("gs://testdataflow123456789/pubsub_to_bq/dataflow/tmp"))
        		//.apply("WriteCounts", TextIO.write().to("gs://testdataflow123456789/pubsub_to_bq/dataflow/tmp"))
                //.via((String json) -> mapJsonToStationUpdate(json))
                //
                // Write the station updates to BigQuery.
                //	
                //.apply(ParDo.of(fn))
        		
                t1.apply("Write to BigQuery", BigQueryIO.<Tweet>write()
                        .to(new TableDestinationFormatter())
                        .withCreateDisposition(CREATE_NEVER)
                        .withWriteDisposition(WRITE_APPEND)
                        .withFormatFunction(new TweetToTableRowMapper()));
                
                // window of 1 minute
                PCollection<Tweet> windowed = t1.apply("Windowing",Window.<Tweet>into(FixedWindows.of(Duration.standardMinutes(4))).withAllowedLateness(Duration.standardMinutes(5)).accumulatingFiredPanes());
                
                PCollection<String> words = windowed.apply("Extract Words", ParDo.of(new DoFn<Tweet, String>() {
                    @ProcessElement
                    public void processElement(ProcessContext c) {
                        // \p{L} denotes the category of Unicode letters,
                        // so this pattern will match on everything that is not a letter.
                    	/*
                        for (String word : c.element().getText().split("[^\\p{L}]+")) {
                            if (!word.isEmpty()) {
                                c.output(word);
                            }
                        }
                        */
                    	if (c.element().getEntities()!=null){
                    		if (c.element().getEntities().getHashtags()!=null){
                    			for (Tweet.Hashtag hashtag : c.element().getEntities().getHashtags()){
                            		c.output(hashtag.getText());
                            	}
                    		}
                    	}
                    }
                }));
                
                PCollection<KV<String, Long>> wordCounts = words.apply(new WordCount.CountWords());
                
                // Filter out late data.
                PCollection<KV<String, Long>> NotLateWordCounts = wordCounts.apply("FilterLateData", ParDo.of(new DoFn<KV<String, Long>, KV<String, Long>>() {
                	@ProcessElement
                    public void processElement(ProcessContext c) {
                		if(c.pane().getTiming() == Timing.LATE) {
                			// Only log late data.
                			logger.info("Late data discarded : Hashtag of value " + c.element().getKey() + " in count ("+c.element().getValue()+") generated at " + c.timestamp());
                		} else {
                			c.output(c.element());
                		}
                	}
                }));

                //PCollection<KV<String,Long>> wordCounts = words.apply("Count words",Count.<String>perElement());

                /**
                 * Concept #5: Format the results and write to a sharded file partitioned by window, using a
                 * simple ParDo operation. Because there may be failures followed by retries, the
                 * writes must be idempotent, but the details of writing to files is elided here.
                 */
                NotLateWordCounts
                    .apply(MapElements.via(new WordCount.FormatAsTextFn()))
                    .apply(new WriteOneFilePerWindow("gs://testdataflow123456789/final/PSTBQ", 1));
                
                /*
                PCollection<String> words = windowed.apply("Extract Words", ParDo.of(new DoFn<Tweet, String>() {
                    @ProcessElement
                    public void processElement(ProcessContext c) {
                        // \p{L} denotes the category of Unicode letters,
                        // so this pattern will match on everything that is not a letter.
                        for (String word : c.element().getText().split("[^\\p{L}]+")) {
                            if (!word.isEmpty()) {
                                c.output(word);
                            }
                        }
                    }
                }));
                
                PCollection<KV<String,Long>> in = words.apply("Count words",Count.<String>perElement());
				*/
                
        p.run().waitUntilFinish();
    }
    
    static class ComputeWordLengthFn extends DoFn<String, Integer> {
    	  @ProcessElement
    	  public void processElement(ProcessContext c) {
    	    // Get the input element from ProcessContext.
    	    String word = c.element();
    	    // Use ProcessContext.output to emit the output element.
    	    c.output(word.length());
    	  }
    	}
    
    static class JsonMap extends DoFn<String, Tweet> {
  	  @ProcessElement
  	  public void processElement(ProcessContext c) {
  	    // Get the input element from ProcessContext.
  	    String json = c.element();
  	    // Use ProcessContext.output to emit the output element.
  	    try {
			c.output(MAPPER.readValue(json, Tweet.class));
		} catch (JsonParseException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (JsonMappingException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
  	  }
  	}
    
    private static Tweet mapJsonToTweet(String json) {
        try {
            return MAPPER.readValue(json, Tweet.class);
        } catch (IOException e) {
            throw new RuntimeException("Unable to map Tweet from string", e);
        }
    }

    public static class TweetToTableRowMapper implements SerializableFunction<Tweet, TableRow> {
        public TableRow apply(Tweet t) {
            final TableRow tr = new TableRow();
            tr.set("Username", t.getUser().getName());
            tr.set("Text", t.getText());
            tr.set("Date",t.getCreated_at());
            if (t.getEntities() == null) {
            	tr.set("Nb_Hashtags", 0);
            } else if (t.getEntities().getHashtags() == null){
            	tr.set("Nb_Hashtags", 0);
            } else {
            	tr.set("Nb_Hashtags", t.getEntities().getHashtags().size());
            }
            tr.set("Language", t.getLang());
            return tr;
        }
    }

    public static class TableDestinationFormatter implements SerializableFunction<ValueInSingleWindow<Tweet>, TableDestination> {

        private static final String DATASET_NAME = "SeriesDataflow";

        private static final DateTimeFormatter PARTITION_FORMATTER = ISODateTimeFormat.basicDate().withZone(DateTimeZone.UTC);

        public TableDestination apply(ValueInSingleWindow<Tweet> value) {
            // Create table reference with only the dataset and table name + partition (no project).
        	/*
            final Tweet su = value.getValue();
            String contract = su.getUser().getName();
            contract = contract != null ? contract.toLowerCase() : "unknown";
            final String partition = PARTITION_FORMATTER.print(value.getTimestamp());
            */
            return new TableDestination(DATASET_NAME + "." + "GOT", null);
        }
    }
}
