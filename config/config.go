package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

// C struct fields are tightly coupled with the keys in config file.
var C Configurations

var configValidator *validator.Validate = validator.New()

// Configurations struct type declares the base config model.
type Configurations struct {
	APPPort    int    `mapstructure:"APP_PORT" validate:"required"`
	DBHost     string `mapstructure:"DB_HOST" validate:"required"`
	DBUser     string `mapstructure:"DB_USER" validate:"required"`
	DBPassword string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBName     string `mapstructure:"DB_NAME" validate:"required"`
	TestDBName string `mapstructure:"TEST_DB_NAME" validate:"required"`
	DBPort     int    `mapstructure:"DB_PORT" validate:"required"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE" validate:"required"`
}

// Initialize initializes the config file and inject into the Configuration struct fields.
func Initialize(configPath string) {
	glog.V(2).Infoln("Initializing configs...")

	viper.SetConfigType("env")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		glog.Fatalf("error reading config file '%s', %s", configPath, err)
	}

	err := viper.Unmarshal(&C)
	if err != nil {
		glog.Fatalf("unable to decode config file '%s' into struct, %v", configPath, err)
	}

	validate()
	glog.V(2).Infoln("Finished initializing configs.")
}

// validate ensures the config values are provided as per defined validation rules.
func validate() {
	err := configValidator.Struct(C)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			glog.Errorf("missing or invalid config value, %s", e)
		}
		glog.Flush()
		os.Exit(1)
	}
}
