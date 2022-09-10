package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dtinth/gojam/pkg/gojam"
)

func main() {
	// Parse command line arguments
	server := flag.String("server", "localhost:22124", "server to connect to")

	flag.Parse()

	fmt.Println("gojam jamulus client")
	fmt.Println("server:", *server)

	client, err := gojam.NewClient(*server)
	if err != nil {
		panic(err)
	}
	fmt.Println("client:", client)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("exiting")
	client.Close()
}
