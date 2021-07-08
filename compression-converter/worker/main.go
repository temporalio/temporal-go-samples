package main

import (
	"log"

	compressionconverter "github.com/temporalio/samples-go/compression-converter"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		// Set DataConverter here so that workflow and activity inputs/results can
		// be compressed/decompressed as required.
		DataConverter: compressionconverter.NewCompressionConverter(
			converter.GetDefaultDataConverter(),
		),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "compressionconverter", worker.Options{})

	w.RegisterWorkflow(compressionconverter.Workflow)
	w.RegisterActivity(compressionconverter.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
