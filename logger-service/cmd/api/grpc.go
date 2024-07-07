package main

import (
	"context"
	"fmt"
	"github.com/olsh4u/protos/gen/go/logs"
	"github.com/org/olsh4u/logger-service/data"
	"google.golang.org/grpc"
	"log"
	"net"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "filed",
		}
		return res, err
	}

	// return resp
	res := &logs.LogResponse{
		Result: "Successfully logged",
	}
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{
		Models: app.Models,
	})
	log.Printf("gRPC server listening on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
