package main

import (
	"database/sql"
	"github.com/olivere/elastic/v7"
	_ "github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/zuyatna/emoney-microservice/transaction-service/server/account-service/pb"
	"github.com/zuyatna/emoney-microservice/transaction-service/server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

}

func runService(logger *logrus.Logger) error {
	cfg, err := config.LoadConfig("..")
	if err != nil {
		logger.WithError(err).Fatal("could not load config")
	}
	logger.Info("Configuration loaded")

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect PostgreSQL")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.WithError(err).Fatal("failed to close database connection")
		} else {
			logger.Info("Database connection closed")
		}
	}(db)

	esClient, err := elastic.NewClient(elastic.SetURL(cfg.ElasticsearchURL), elastic.SetSniff(false))
	if err != nil {
		logger.WithError(err).Fatal("failed to create Elasticsearch client")
	}

	accountServiceConn, err := grpc.NewClient(cfg.AccountServiceTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to account service")
	}
	defer func(accountServiceConn *grpc.ClientConn) {
		err := accountServiceConn.Close()
		if err != nil {
			logger.WithError(err).Fatal("failed to close account service connection")
		} else {
			logger.Info("Account service connection closed")
		}
	}(accountServiceConn)
	accountServiceClient := pb.NewAccountServiceClient(accountServiceConn)

	rabbitConn, err := amqp.Dial(cfg.RABBITMQURL)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to RabbitMQ")
	}
	defer func(rabbitConn *amqp.Connection) {
		err := rabbitConn.Close()
		if err != nil {
			logger.WithError(err).Fatal("failed to close RabbitMQ connection")
		} else {
			logger.Info("RabbitMQ connection closed")
		}
	}(rabbitConn)
	logger.Info("Connected to RabbitMQ")

	// TODO: Implement the rest of the service logic here, such as setting up gRPC servers, handling requests, etc.

	errChan := make(chan error, 2)
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			errChan <- err
			return
		}
		logger.WithField("port", cfg.GRPCPort).Info("gRPC server listening")

	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

}
