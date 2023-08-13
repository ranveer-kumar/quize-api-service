package main

// import (
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// )

// const version = "1.0.0"

// type config struct {
// 	port int
// 	env string
// }

// type application struct {
// 	config config
// 	logger *log.Logger

// }

// func main() {
// 	var cfg config
// 	flag.IntVar(&cfg.port, "port", 4000, "API server port")
// 	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
// 	flag.Parse()

// 	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

// 	app := &application{
//         config: cfg,
//         logger: logger,
//     }
//
// }

import (
	"log"
	"net/http"
	"quize-api-service/configs"
	"quize-api-service/internal/db"
	"quize-api-service/internal/logger"
	"quize-api-service/internal/routes"
)

func main() {

	logger.InitLogger()

	cfg := configs.LoadConfig() // Load your application configuration

	// Connect to MongoDB
	db.InitMongoDB(cfg.MongoURI, cfg.MongoDBName)

	router := routes.SetupRouter()

	// Start the server
	serverAddr := ":" + cfg.ServerPort
	logger.Info("Server started at " + serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
