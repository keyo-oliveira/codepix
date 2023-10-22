package cmd

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/keyo-oliveira/codepix/application/grpc"
	"github.com/keyo-oliveira/codepix/application/kafka"
	"github.com/keyo-oliveira/codepix/infra/db"
	"github.com/spf13/cobra"
	"os"
)

var gRPCPortNumber int

// kafkaCmd represents the kafka command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Start all servers ( grpc and kafka consumer )",

	Run: func(cmd *cobra.Command, args []string) {
		database = db.ConnectDB(os.Getenv("env"))
		go grpc.StartGrpcServer(database, gRPCPortNumber)

		deliveryChan := make(chan ckafka.Event)
		producer := kafka.TheNewKafkaProducer()
		kafka.Publish("Ola kafka", "teste", producer, deliveryChan)
		go kafka.DeliveryReport(deliveryChan)

		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		kafkaProcessor.Consume()
		print("KAFKA AND GRPC SERVER RUNNING ON")

	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	grpcCmd.Flags().IntVarP(&gRPCPortNumber, "grpc-port", "q", 50051, "gRPC Server Port")
}
