package main

import (
	"context"
	"log"
	"os"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
