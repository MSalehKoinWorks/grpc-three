package account

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/MSalehKoinWorks/grpc-three/proto/account"
)

type DepositClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
}

func NewDepositClient(conn *grpc.ClientConn, timeout time.Duration) DepositClient {
	return DepositClient{
		conn:    conn,
		timeout: timeout,
	}
}

func (d DepositClient) Deposit(ctx context.Context, amount float32) (bool, error) {
	client := pb.NewDepositServiceClient(d.conn)
	request := &pb.Req{Amount: amount}
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(d.timeout))

	defer cancel()

	response, err := client.Deposit(ctx, request)

	if err != nil {
		if er, ok := status.FromError(err); ok {
			return false, fmt.Errorf("grpc: %s, %s", er.Code(), er.Message())
		}
		return false, fmt.Errorf("server: %s", err.Error())
	}

	return response.GetOk(), nil
}

func (d DepositClient) Withdraw(ctx context.Context, amount float32) (bool, error) {
	client := pb.NewDepositServiceClient(d.conn)
	request := &pb.Req{Amount: amount}
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(d.timeout))

	defer cancel()

	response, err := client.Withdraw(ctx, request)

	if err != nil {
		if er, ok := status.FromError(err); ok {
			return false, fmt.Errorf("grpc: %s, %s", er.Code(), er.Message())
		}
		return false, fmt.Errorf("server: %s", err.Error())
	}

	return response.GetOk(), nil
}
