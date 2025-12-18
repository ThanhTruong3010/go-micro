package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/rpc"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) RPCListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	// 0.0.0.0 listening on all the configured network interfaces
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database(("logs")).Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	})

	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Name

	return nil

}
