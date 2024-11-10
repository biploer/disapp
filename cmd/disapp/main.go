package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"disapp/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
