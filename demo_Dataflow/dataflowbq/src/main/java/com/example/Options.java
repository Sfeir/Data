package com.example;

import org.apache.beam.sdk.options.Description;
import org.apache.beam.sdk.options.PipelineOptions;
import org.apache.beam.sdk.options.Validation.Required;
import org.apache.beam.sdk.options.ValueProvider;

public interface Options extends PipelineOptions {

    @Description("Pub/Sub subscription from which to read the station updates")
    @Required
    ValueProvider<String> getSubscription();
    void setSubscription(ValueProvider<String> value);

    @Description("Table identifier (<project>:<dataset>.<table_name>) on which to write the station updates")
    @Required
    ValueProvider<String> getTable();
    void setTable(ValueProvider<String> value);
}
