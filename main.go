package main

import (
	"log"
	"omni-cli/app"
	"os"
)

func main() {

	aplicacao := app.Build()
	if erro := aplicacao.Run(os.Args); erro != nil {
		log.Fatal(erro)
	}
}
