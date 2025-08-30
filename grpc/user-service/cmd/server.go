package main

import (
	"context"
	"user-service/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.NewApp(ctx)

	if err := app.Run(); err != nil {
		panic(err)
	}

	defer app.Shutdown()
}
