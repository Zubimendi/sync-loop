package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"github.com/Zubimendi/sync-loop/api/internal/activity"
	"github.com/Zubimendi/sync-loop/api/internal/workflow"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: os.Getenv("TEMPORAL_HOST"),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "sync-loop-task-queue", worker.Options{})
	w.RegisterWorkflow(workflow.CopyTableWorkflow)
	w.RegisterActivity(activity.CopyTableActivity)

	log.Println("Worker started")
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}