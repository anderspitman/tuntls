package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {

	server := flag.String("server", "", "TLS Server")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "tuntls connecting to server: %s\n", *server)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", *server), &tls.Config{
		//RootCAs: roots,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: "+err.Error())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(conn, os.Stdin)
		wg.Done()
	}()

	go func() {
		io.Copy(os.Stdout, conn)
		wg.Done()
	}()

	wg.Wait()

	conn.Close()
}
