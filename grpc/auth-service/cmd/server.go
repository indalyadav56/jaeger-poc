package main

import (
	"auth-service/internal/app"
	"context"
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
