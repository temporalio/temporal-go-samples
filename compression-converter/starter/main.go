package main

import (
	"context"
	"log"

	compressionconverter "github.com/temporalio/samples-go/compression-converter"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		// Set DataConverter here to ensure that workflow inputs and results are
		// compressed/decompressed as required.
		DataConverter: compressionconverter.NewCompressionConverter(
			converter.GetDefaultDataConverter(),
		),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "compressionconverter_workflowID",
		TaskQueue: "compressionconverter",
	}

	ctx := context.Background()

	// The workflow input "Hello" will be compressed by the DataConverter before being sent to Temporal
	we, err := c.ExecuteWorkflow(
		ctx,
		workflowOptions,
		compressionconverter.Workflow,
		"My Friend",
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
