package main

import (
	"context"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/secrets"
)

func main() {
	log := logger.All()
	client := secrets.Must(secrets.New(secrets.OptConfigFromEnv(), secrets.OptLog(log)))

	key := "cubbyhole/willtest"

	ctx := context.Background()

	if err := client.Put(ctx, key, secrets.Values{"value": "THE FOOOS"}); err != nil {
		log.Fatal(err)
		return
	}
	if err := client.Put(ctx, key, secrets.Values{"value": "THE BUZZ"}); err != nil {
		log.Fatal(err)
		return
	}

	values, err := client.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("values: %#v", values)

	if err := client.Delete(ctx, key); err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("~fin~")
}
