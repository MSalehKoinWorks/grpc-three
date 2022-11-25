package account

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/MSalehKoinWorks/grpc-three/proto/account"
)

var amount float32 = 1000

type mockDepositServer struct {
	pb.UnimplementedDepositServiceServer
}

func (*mockDepositServer) Deposit(ctx context.Context, req *pb.Req) (*pb.Res, error) {
	if req.GetAmount() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot deposit %v", req.GetAmount())
	}

	amount += req.GetAmount()
	fmt.Printf("deposit: %v, new balance: %v\n", req.GetAmount(), amount)

	return &pb.Res{Ok: true}, nil
}

func (*mockDepositServer) Withdraw(ctx context.Context, req *pb.Req) (*pb.Res, error) {
	if req.GetAmount() < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "cannot withdraw %v", req.GetAmount())
	}

	if req.GetAmount() > amount {
		return nil, status.Errorf(codes.InvalidArgument, "cannot withdraw %v", req.GetAmount())
	}

	amount -= req.GetAmount()
	fmt.Printf("withdraw: %v, new balance: %v\n", req.GetAmount(), amount)

	return &pb.Res{Ok: true}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	pb.RegisterDepositServiceServer(server, &mockDepositServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestDepositServer_Deposit(t *testing.T) {
	tests := []struct {
		name    string
		amount  float32
		res     *pb.Res
		errCode codes.Code
		errMsg  string
	}{
		{
			"invalid request with negative amount",
			-100000.00,
			nil,
			codes.InvalidArgument,
			fmt.Sprintf("cannot deposit %v", -100000.00),
		},
		{
			"valid request with non negative amount",
			100000.00,
			&pb.Res{Ok: true},
			codes.OK,
			"",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDepositServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pb.Req{Amount: tt.amount}

			response, err := client.Deposit(ctx, request)

			if response != nil {
				if response.GetOk() != tt.res.GetOk() {
					t.Error("response: expected", tt.res.GetOk(), "received", response.GetOk())
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}

func TestDepositServer_Withdraw(t *testing.T) {
	tests := []struct {
		name    string
		amount  float32
		res     *pb.Res
		errCode codes.Code
		errMsg  string
	}{
		{
			"invalid request with negative amount",
			-500.00,
			nil,
			codes.InvalidArgument,
			fmt.Sprintf("cannot withdraw %v", -500.00),
		},
		{
			"valid request with non negative amount",
			500.00,
			&pb.Res{Ok: true},
			codes.OK,
			"",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDepositServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pb.Req{Amount: tt.amount}

			response, err := client.Withdraw(ctx, request)

			if response != nil {
				if response.GetOk() != tt.res.GetOk() {
					t.Error("response: expected", tt.res.GetOk(), "received", response.GetOk())
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}
