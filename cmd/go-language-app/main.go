package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MihirMohapatra/go-cli-app/internal/greeting"
)

func main() {
	name := flag.String("name", "World", "name to greet")
	flag.Parse()

	message, err := greeting.Build(*name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(message)
}
