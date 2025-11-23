package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GCP      GCPConfig
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release
}

type DatabaseConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type GCPConfig struct {
	ProjectID           string `mapstructure:"project_id"`
	StorageBucket       string `mapstructure:"storage_bucket"`
	CredentialsJSONPath string `mapstructure:"credentials_json_path"`
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.database", "taiwanstay")

	// Bind environment variables
	// Example: SERVER_PORT maps to Server.Port
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.mode", "GIN_MODE")
	_ = viper.BindEnv("database.uri", "MONGODB_URI")
	_ = viper.BindEnv("database.database", "MONGODB_DATABASE")
	_ = viper.BindEnv("gcp.project_id", "GCP_PROJECT_ID")
	_ = viper.BindEnv("gcp.storage_bucket", "GCP_STORAGE_BUCKET")
	_ = viper.BindEnv("gcp.credentials_json_path", "GOOGLE_APPLICATION_CREDENTIALS")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
