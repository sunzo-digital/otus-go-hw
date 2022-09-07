package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient struct {
	network, address string
	timeout          time.Duration
	in               io.ReadCloser
	out              io.Writer
	connection       net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *TelnetClient {
	return &TelnetClient{
		network: "tcp",
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (client *TelnetClient) Connect() error {
	connection, err := net.DialTimeout(client.network, client.address, client.timeout)
	if err != nil {
		return err
	}

	client.connection = connection
	_, _ = fmt.Fprintf(os.Stderr, "...Connected to %s\n", client.address)
	return nil
}

func (client *TelnetClient) Close() error {
	return client.connection.Close()
}

func (client *TelnetClient) Send() error {
	_, err := io.Copy(client.connection, client.in)
	return err
}

func (client *TelnetClient) Receive() error {
	_, err := io.Copy(client.out, client.connection)
	return err
}
