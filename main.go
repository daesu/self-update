package main

import (
	"log"
	"self-update/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}
