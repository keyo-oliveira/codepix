package cmd

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/keyo-oliveira/codepix/application/grpc"
	"github.com/keyo-oliveira/codepix/infra/db"
	"os"

	"github.com/spf13/cobra"
)

var portName int
var database *gorm.DB

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start a new grpc server",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grpc called")
		database = db.ConnectDB(os.Getenv("env"))
		grpc.StartGrpcServer(database, portName)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().IntVarP(&portName, "port", "p", 50051, "gRPC Server Port")
}
