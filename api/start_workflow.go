package main

import (
	"context"
	"log"
	"go.temporal.io/sdk/client"
	"github.com/Zubimendi/sync-loop/api/internal/workflow"
	"time"
)

func main() {
	c, err := client.Dial(client.Options{HostPort: "localhost:7233"})
	if err != nil { log.Fatal(err) }
	defer c.Close()

	we, err := c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: "sync-loop-task-queue",
			ID:        "copy-users-" + time.Now().Format("20060102-150405"),
		},
		workflow.CopyTableWorkflow,
		"users",
	)
	if err != nil { log.Fatal(err) }
	log.Printf("started workflow %s %s", we.GetID(), we.GetRunID())
}