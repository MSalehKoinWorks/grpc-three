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

	for i := 1; i <= 25; i++ {
		time.Sleep(time.Millisecond * 10)
		doDeposit(conn, float32(i)*1000)
	}

	for i := 1; i <= 25; i++ {
		time.Sleep(time.Millisecond * 10)
		doWithdraw(conn, float32(i)*500)
	}
}

func doDeposit(conn *grpc.ClientConn, amount float32) {
	response, err := account.NewDepositClient(conn, time.Second).Deposit(context.Background(), amount)

	if err != nil {
		log.Fatalln(err)
	}

	if !response {
		log.Println("Deposit failed ...")
	}

	log.Println("Deposit success ...")
}

func doWithdraw(conn *grpc.ClientConn, amount float32) {
	response, err := account.NewDepositClient(conn, time.Second).Withdraw(context.Background(), amount)

	if err != nil {
		log.Fatalln(err)
	}

	if !response {
		log.Println("Withdraw failed ...")
	}

	log.Println("Withdraw success ...")
}
