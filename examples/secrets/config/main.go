package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/secrets"
	"github.com/blend/go-sdk/yaml"
)

type myConfig struct {
	Environment string `yaml:"environment" env:"SERVICE_ENV" secret:"environment"`
	Version     string `yaml:"version" env:"SERVICE_VERSION"`
	Secret      string `yaml:"secret" secret:"secret,base64"`
}

func main() {
	log := logger.All()
	client := secrets.Must(secrets.New(secrets.OptConfigFromEnv(), secrets.OptLog(log)))

	keyPath := "configTest"
	err := client.WriteInto(context.TODO(), keyPath, myConfig{
		Environment: "test",
		Version:     "wont be in output",
		Secret:      base64.StdEncoding.EncodeToString([]byte("a super secure one")),
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cfg myConfig
	err = client.ReadInto(context.TODO(), keyPath, &cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	contents, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(contents))
}
