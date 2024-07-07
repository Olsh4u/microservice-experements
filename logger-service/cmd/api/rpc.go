package main

import (
	"context"
	"github.com/org/olsh4u/logger-service/data"
	"log"
	"time"
)

// RPCServer is the type for our RPC Server. Method that takes this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct {
}

// RPCPayload is the type for data we receive from RPC
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo writes our payload to mongo
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log into Mongo DB", err)
	}

	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
