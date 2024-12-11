package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"
	"time"
)

// we have a configuration
// we get a request
// check if there is an entry for this request
// check rewriteRules
// we send it to that server

// Lets have 100 requests running that all of them make some requests through the proxy
// lets use a custom transport network call

type ProxyConfiguration struct {
	in  string
	out string
}

var proxyConfigs = []ProxyConfiguration{
	{
		in:  "/api/users",
		out: "localhost:8082",
	},
	{
		in:  "/api/products",
		out: "localhost:8081",
	},
	{
		in:  "/api/orders",
		out: "localhost:8081",
	},
	{
		in:  "/api/auth",
		out: "localhost:8081",
	},
	{
		in:  "/api/payments",
		out: "localhost:8081",
	},
	{
		in:  "/api/inventory",
		out: "localhost:8081",
	},
	{
		in:  "/api/shipping",
		out: "localhost:8081",
	},
	{
		in:  "/api/notifications",
		out: "localhost:8081",
	},
	{
		in:  "/api/analytics",
		out: "localhost:8081",
	},
	{
		in:  "/api/admin",
		out: "localhost:8081",
	},
}

func dialUp(host string, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", host)
	if err != nil {
		log.Printf("Failed to dial: %v", err)
		return
	}
	if _, err := conn.Write([]byte(host + " " + message)); err != nil {
		log.Fatal(err)
	}
}

func validateEelProtocol(buffer string) (string, string, error) {

	transmitted_buffer := strings.Split(buffer, " ")
	if transmitted_buffer[0] != "eel" && len(transmitted_buffer) < 3 {
		return "", "", errors.New("not a properly formed Eel protocol")

	}

	return transmitted_buffer[1], strings.Join(transmitted_buffer[2:], " "), nil
}

func proxy(conn net.Conn) {
	defer conn.Close()

	for {
		var buffer [124]byte
		n, err := conn.Read(buffer[:])

		if err != nil {
			return
		}

		transmitted_buffer := string(buffer[:n])
		uri, message, err := validateEelProtocol(transmitted_buffer)
		if err != nil {
			fmt.Printf("Invalid protocol data format: %v", transmitted_buffer)
			return
		}

		idx := slices.IndexFunc(
			proxyConfigs,
			func(c ProxyConfiguration) bool { return c.in == uri },
		)
		if idx == -1 {
			return
		}

		dialUp(proxyConfigs[idx].out, message)
	}
}

func main() {
	port := ":3000"
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("couldn't listen to network")
	}
	defer l.Close()
	fmt.Println("Proxy server up and running at " + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln("err while accepting", err)
			continue
		}
		go proxy(conn)

	}
}
