package main

import (
	"context"
	"go-nsq/application/entrypoint"
	"go-nsq/application/mq/nsq"
	"go-nsq/application/mq/redis"
	"go-nsq/db"
	"go-nsq/store/minio"
	"go-nsq/store/nosql"
	"go-nsq/transport"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := db.InitMongoDB(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer client.Client.Disconnect(ctx)

	mongoDBStore := nosql.NewNoSQLStore(client)

	nsq := nsq.NewNSQClient()
	minio, err := minio.InitMinioService(ctx, "documents")
	if err != nil {
		log.Fatalf("Error intialize Minio Client")
	}
	redis, err := redis.NewRedisClient()
	entryPointService := entrypoint.NewEntryPointService(mongoDBStore, nsq, minio, redis)
	server := transport.NewHTTPServer(entryPointService)
	serverAddr := os.Getenv("SERVER_ADDR")
	err = http.ListenAndServe(serverAddr, server)
	if err != nil {
		log.Fatalf("Error connect to the %s port \n", serverAddr)
	}
}
