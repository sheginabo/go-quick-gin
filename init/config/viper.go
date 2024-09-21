package config

import (
	"github.com/spf13/viper"
	"log"
)

func NewModule(path string) {
	// set default values
	setDefault()

	// 這裡的 "./" 代表目前工作目錄，也就是你在命令列中執行 Go 指令時所處的目錄。
	// "./" here represents the current working directory, which is the directory where you run the Go command from in the command line.
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// check if ENV is set
	env := viper.Get("ENV")
	if env == "" {
		log.Fatalf("Environment variable ENV not set or empty")
	}
	log.Printf("Environment variable ENV: %s", env)
}

func setDefault() {
	viper.SetDefault("ENV", "local")
	viper.SetDefault("APP_NAME", "go-quick-gin")
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")
}
