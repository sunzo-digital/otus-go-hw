package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	client, err := initClient()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer func(client *TelnetClient) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

	runWorkers(client)
}

func initClient() (client *TelnetClient, err error) {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		return nil, errors.New("invalid arguments count. 2 were expected")
	}

	return NewTelnetClient(
		net.JoinHostPort(args[0], args[1]),
		timeout,
		os.Stdin,
		os.Stdout,
	), nil
}

func runWorkers(client *TelnetClient) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := client.Send(); err != nil {
			log.Println("Sending error: ", err)
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			log.Println("Receiving error: ", err)
		}
		_, _ = fmt.Fprint(os.Stderr, "...Connection was closed by peer\n")
	}()

	<-ctx.Done()
}
