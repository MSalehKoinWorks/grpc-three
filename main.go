package main

import (
	"context"
	"log"
	"time"

	"github.com/MSalehKoinWorks/grpc-three/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Client running ...")

	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	response, err := client.
		NewDepositClient(conn, time.Second).
		Deposit(context.Background(), 1990.01)

	log.Println(response)
	log.Println(err)
}
