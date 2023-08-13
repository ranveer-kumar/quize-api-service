package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string
	MongoURI    string
	MongoDBName string
	// Add other configuration parameters
}

func LoadConfig() *Config {
	viper.SetConfigName("quize-config")
	viper.AddConfigPath(".")    // You can change this to the path of your config file
	viper.SetConfigType("yaml") // or "json", "toml", etc.

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	return &cfg
}
