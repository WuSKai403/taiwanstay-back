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
	Image    ImageConfig
}

type ServerConfig struct {
	Port      string `mapstructure:"port"`
	Mode      string `mapstructure:"mode"` // debug, release
	JWTSecret string `mapstructure:"jwt_secret"`
}

type DatabaseConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type GCPConfig struct {
	ProjectID           string `mapstructure:"project_id"`
	PublicBucket        string `mapstructure:"public_bucket"`
	PrivateBucket       string `mapstructure:"private_bucket"`
	CredentialsJSONPath string `mapstructure:"credentials_json_path"`
}

type ImageConfig struct {
	ImageKitEndpoint string `mapstructure:"imagekit_endpoint"`

	// Reject Thresholds (>= this value -> REJECT)
	RejectAdult    string `mapstructure:"reject_adult"`
	RejectSpoof    string `mapstructure:"reject_spoof"`
	RejectMedical  string `mapstructure:"reject_medical"`
	RejectViolence string `mapstructure:"reject_violence"`
	RejectRacy     string `mapstructure:"reject_racy"`
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

	// Image Config Defaults
	viper.SetDefault("image.imagekit_endpoint", "")

	// Default Reject Thresholds (Strict on Violence/Adult)
	viper.SetDefault("image.reject_adult", "LIKELY")
	viper.SetDefault("image.reject_spoof", "LIKELY")
	viper.SetDefault("image.reject_medical", "LIKELY")
	viper.SetDefault("image.reject_violence", "LIKELY")
	viper.SetDefault("image.reject_racy", "LIKELY")

	// Bind environment variables
	// Example: SERVER_PORT maps to Server.Port
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.mode", "GIN_MODE")
	_ = viper.BindEnv("server.jwt_secret", "JWT_SECRET")
	_ = viper.BindEnv("database.uri", "MONGODB_URI")
	_ = viper.BindEnv("database.database", "MONGODB_DATABASE")
	_ = viper.BindEnv("gcp.project_id", "GCP_PROJECT_ID")
	_ = viper.BindEnv("gcp.public_bucket", "GCP_STORAGE_PUBLIC_BUCKET")
	_ = viper.BindEnv("gcp.private_bucket", "GCP_STORAGE_PRIVATE_BUCKET")
	_ = viper.BindEnv("gcp.credentials_json_path", "GOOGLE_APPLICATION_CREDENTIALS")

	_ = viper.BindEnv("image.imagekit_endpoint", "IMAGEKIT_URL_ENDPOINT")

	_ = viper.BindEnv("image.reject_adult", "IMAGE_REJECT_ADULT")
	_ = viper.BindEnv("image.reject_spoof", "IMAGE_REJECT_SPOOF")
	_ = viper.BindEnv("image.reject_medical", "IMAGE_REJECT_MEDICAL")
	_ = viper.BindEnv("image.reject_violence", "IMAGE_REJECT_VIOLENCE")
	_ = viper.BindEnv("image.reject_racy", "IMAGE_REJECT_RACY")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
