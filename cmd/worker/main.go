package main

import (
	"doctravel/internal/docactivity"
	"doctravel/internal/docflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	// Initialize a Temporal Client
	// Specify the Namespace in the Client options
	clientOptions := client.Options{
		Namespace: docflow.Namespace,
	}
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create a Temporal Client", err)
	}
	defer temporalClient.Close()

	// Create a new Worker
	yourWorker := worker.New(temporalClient, "doc_travel", worker.Options{})
	// Register Workflows
	yourWorker.RegisterWorkflow(docflow.DocProcessingFlow)
	// Register Activities
	yourWorker.RegisterActivity(docactivity.DocAction{}.ApproveActivity)
	yourWorker.RegisterActivity(docactivity.DocAction{}.StartEditingActivity)
	// Start the Worker Process
	err = yourWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start the Worker Process", err)
	}
}
