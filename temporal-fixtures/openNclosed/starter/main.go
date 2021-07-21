package main

import (
	"context"
	// "fmt"
	"log"

	"github.com/pborman/uuid"
	"github.com/temporalio/samples-go/temporal-fixtures/openNclosed"

	"strconv"

	"go.temporal.io/sdk/client"
)

var (
	NumberOfWorkflows = 5
	KeepOpen          = false
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{Namespace: "default"})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	uuidvar := uuid.New()
	i := 1
	for i <= NumberOfWorkflows {
		id := "open-n-closed##" + uuidvar[:6] + "##" + strconv.Itoa(i)
		i++

		workflowOptions := client.StartWorkflowOptions{
			ID:        id,
			TaskQueue: "open-n-closed",
		}

		_, err := c.ExecuteWorkflow(context.Background(), workflowOptions, openNclosed.OpenClosedFixtureWorkflow,
			"Temporal", KeepOpen)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
	}
}
