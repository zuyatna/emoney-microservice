package handler

import "github.com/zuyatna/emoney-microservice/transaction-service/server/pb"

type TransactionHandler struct {
	pb.UnimplementedTransactionServiceServer
	
}