package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/zuyatna/emoney-microservice/account-service/server/config"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	app := &cli.App{
		Name:  "account-service",
		Usage: "A microservice for managing user accounts",
		Action: func(c *cli.Context) error {
			return runService(logger)
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatalf("Failedto run CLI app: %v", err)
	}
}

func runService(logger *logrus.Logger) error {
	cfg, err := config.LoadConfig("..")
	if err != nil {
		logger.Fatalf("Error loading configuration: %v", err)
	}

	ctx := context.Background()
	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		logger.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatalf("Error closing db: %v", err)
		}
	}(db)

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err := db.PingContext(ctx); err != nil {
		logger.Fatalf("Error pinging PostgreSQL: %v", err)
	}
	logger.Println("Connected to PostgreSQL")

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Fatalf("Error parsing Redis URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.Fatalf("Error connecting to Redis: %v", err)
	}
	logger.Println("Connected to Redis")

	rabbitConn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer func(rabbitConn *amqp.Connection) {
		err := rabbitConn.Close()
		if err != nil {
			logger.Fatalf("Error closing connection: %v", err)
		}
	}(rabbitConn)
	logger.Println("Connected to RabbitMQ")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down gracefully...")

	return nil
}
