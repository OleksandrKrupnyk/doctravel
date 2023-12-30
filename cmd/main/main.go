package main

import (
	"context"
	"doctravel/internal/docactivity"
	"doctravel/internal/docflow"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"log"
	"log/slog"
)

const (
	ID = "doc_travel_workflowID"
)

func main() {

	// Initialize a Temporal Client
	// Specify the IP, port, and Namespace in the Client options
	options := client.Options{
		HostPort:  "127.0.0.1:7233",
		Namespace: docflow.Namespace,
	}
	clientOptions := options
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create a Temporal Client", err)
	}
	defer temporalClient.Close()

	app := fiber.New()
	app.Post("/new", func(c *fiber.Ctx) error {
		workflowOptions := client.StartWorkflowOptions{
			ID:        ID + uuid.NewUUID().String(), // ID виконання
			TaskQueue: "doc_travel",                 // Черга в яку поставити виконання
		}
		docA := docactivity.DocAction{}
		// Запуск потоку виконання
		we, err := temporalClient.ExecuteWorkflow(context.Background(),
			workflowOptions, docflow.DocProcessingFlow,
			docA,
			"Test",
		)
		if err != nil {
			log.Println("Unable to execute workflow", err)
		}
		s := fmt.Sprintf("Started workflow WorkflowID %s  RunID %s", we.GetID(), we.GetRunID())
		slog.Info(s)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"id":     we.GetID(),
			"run_id": we.GetRunID(),
		})
	})

	app.Put("/edit/:id/:run_id", func(c *fiber.Ctx) error {
		runID := c.Params("run_id", "")
		id := c.Params("id", "")
		signal := docflow.EditSignal{
			Editing: "Toster",
		}
		err = temporalClient.SignalWorkflow(context.Background(), id, runID, docflow.EditingSignalName, signal)
		if err != nil {
			log.Println("Error sending the EditingSignalName", err)
			return c.Status(fiber.StatusBadRequest).SendString("Error sending the EditingSignalName")
		}
		return c.SendString("EditSignal Ok")
	})

	app.Put("/approve/:id/:run_id", func(c *fiber.Ctx) error {
		runID := c.Params("run_id", "")
		id := c.Params("id", "")
		signal := docflow.ApproveSignal{
			Approved: "Admin",
		}
		err = temporalClient.SignalWorkflow(context.Background(), id, runID, docflow.ApproveSignalName, signal)
		if err != nil {
			log.Println("Error sending the ApproveSignalName", err)
			return c.Status(fiber.StatusBadRequest).SendString("Error sending the ApproveSignalName")
		}
		return c.SendString("ApproveSignal Ok")
	})

	app.Get("/:id/:run_id", func(c *fiber.Ctx) error {
		runID := c.Params("run_id", "")
		id := c.Params("id", "")
		log.Println("runID", runID)
		response, err := temporalClient.QueryWorkflow(context.Background(), id, runID, docflow.InfoQuery, "[", "]")
		if err != nil {
			log.Println("Unable to get Query to client: ", err)
			return c.Status(fiber.StatusBadRequest).SendString("Unable to get Query to client: " + err.Error())
		}
		var res string
		if err := response.Get(&res); err != nil {
			log.Println("Error convert Result: ", err)
			return c.Status(fiber.StatusBadRequest).SendString("Error convert Result: " + err.Error())
		}
		return c.SendString(res)
	})

	app.Listen(":3000")

}
