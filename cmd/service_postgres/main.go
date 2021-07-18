package main

import (
	"github.com/hywmongous/example-service/internal/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.Module).Run()
}
