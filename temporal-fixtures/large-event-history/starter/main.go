package main

import (
	"context"
	"log"

	"github.com/pborman/uuid"
	largeeventhistory "github.com/temporalio/samples-go/temporal-fixtures/large-event-history"
	"go.temporal.io/sdk/client"
)

var (
	LengthOfHistory = 1000
	WillFailOrNot   = true
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	id := uuid.New()[0:4]
	workflowOptions := client.StartWorkflowOptions{
		ID:        "largeeventhistory_" + id,
		TaskQueue: "largeeventhistory",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions,
		largeeventhistory.LargeEventHistoryWorkflow, LengthOfHistory, WillFailOrNot)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
