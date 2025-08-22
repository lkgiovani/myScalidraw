package main

import (
	"go.uber.org/fx"

	fxModules "myScalidraw/infra/fx"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load(".env")

	fx.New(
		fxModules.AllModules,
	).Run()
}
