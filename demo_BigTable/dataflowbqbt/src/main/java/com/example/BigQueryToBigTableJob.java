package com.example;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import org.apache.beam.sdk.Pipeline;
import org.apache.beam.sdk.io.gcp.bigquery.BigQueryIO;
import org.apache.beam.sdk.io.gcp.bigtable.BigtableIO;
import org.apache.beam.sdk.options.PipelineOptionsFactory;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.values.KV;
import org.apache.beam.sdk.values.PCollection;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.api.services.bigquery.model.TableReference;
import com.google.api.services.bigquery.model.TableRow;
import com.google.bigtable.v2.Mutation;
import com.google.cloud.bigtable.config.BigtableOptions;
import com.google.protobuf.ByteString;

public class BigQueryToBigTableJob {

    private static final Logger logger = LoggerFactory.getLogger(BigQueryToBigTableJob.class);

    /*
    public static CloudBigtableScanConfiguration config = new CloudBigtableScanConfiguration.Builder()
    	    .withProjectId("sfeir-data")
    	    .withInstanceId("tweets")
    	    .withTableId("got")
    	    .build();
    */
    public static void main(String[] args) {

    	 BigtableOptions.Builder optionsBuilder =
    			      new BigtableOptions.Builder()
    			          .setProjectId("sfeir-data")
    			          .setInstanceId("tweets");
    			 
    			  //PCollection<KV<ByteString, Iterable<Mutation>>> data = ...;
    			 
    	
        // Get the job's options from the command line arguments.
        final Options options = PipelineOptionsFactory.fromArgs(args).withValidation().as(Options.class);

        // Create the pipeline.
        final Pipeline p = Pipeline.create(options);
        TableReference table = new TableReference();
        table.setProjectId("sfeir-data");
        table.setDatasetId("SeriesDataflow");
        table.setTableId("GOT");
        PCollection<TableRow> tablerow = p.apply("Read from BigQuery",
        //BigQueryIO.read().from(table));
		//BigQueryIO.read().fromQuery("select Username,Text,Date,Language from [SeriesDataflow.GOT] where Date is not null"));
        BigQueryIO.read().fromQuery("select Username,Text, FORMAT_DATETIME('%Y%m%d%H%M%S',PARSE_DATETIME('%a %b %d %X +0000 %Y',Date)) AS Date from `SeriesDataflow.GOT` where Date is not null and length(Date)=30 and length(Username)>2").usingStandardSql());
        PCollection<KV<ByteString, Iterable<Mutation>>> data = tablerow.apply(ParDo.of(MUTATION_TRANSFORM));
		data.apply(BigtableIO.write().withBigtableOptions(optionsBuilder).withTableId("got"));
		p.run().waitUntilFinish();
	}
    


  static final DoFn<TableRow, KV<ByteString, Iterable<Mutation>>> MUTATION_TRANSFORM = new DoFn<TableRow, KV<ByteString, Iterable<Mutation>>>() {
    private static final long serialVersionUID = 1L;

    @ProcessElement
    public void processElement(DoFn<TableRow, KV<ByteString, Iterable<Mutation>>>.ProcessContext c) throws Exception {

      TableRow row = c.element();

      List<Mutation> mutations = new ArrayList<>();
      
      for (Map.Entry<String, Object> field : row.entrySet()) {
    	    			              Mutation m = Mutation.newBuilder()
    			                      .setSetCell(
    			                          Mutation.SetCell.newBuilder()
    			                              .setValue(ByteString.copyFromUtf8(String.valueOf(field.getValue())))
    			                              .setFamilyName("data")
    			                              .setColumnQualifier(ByteString.copyFromUtf8((field.getKey()))))
    			                      .build();
    	    			              mutations.add(m);
    			              
      }
      //ByteString key = ByteString.copyFromUtf8(String.valueOf("key"));
      //ByteString key = ByteString.copyFromUtf8(UUID.randomUUID().toString());
      ByteString key = ByteString.copyFromUtf8(row.get("Username")+"#"+row.get("Date"));
      c.output(KV.of(key, mutations));
    }
  };
  
}
