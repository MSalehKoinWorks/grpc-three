package main

import (
	"context"
	"log"
	"time"

	"github.com/MSalehKoinWorks/grpc-three/controllers/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Client running ...")

	conn, err := grpc.Dial(":6666", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	response, err := account.NewDepositClient(conn, time.Second).Deposit(context.Background(), 1000)

	if err != nil {
		log.Fatalln(err)
	}

	if !response {
		log.Println("Deposit failed ...")
	}

	log.Println("Deposit success ...")
}
