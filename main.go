package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func main() {

	server := flag.String("server", "", "TLS Server")
	port := flag.Int("port", 0, "Local port to bind to")
	flag.Parse()

	if *port == 0 {
		// one-time tunnel over stdin/stdout
		doTunnel(*server, os.Stdin, os.Stdout)
	} else {
		// listen on a port and create tunnels for each connection
		fmt.Fprintf(os.Stderr, "Listening on port %d\n", *port)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(1)
			}

			go doTunnel(*server, conn, conn)
		}
	}

}

func doTunnel(server string, in io.Reader, out io.Writer) {
	fmt.Fprintf(os.Stderr, "tuntls connecting to server: %s\n", server)

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", server), &tls.Config{
		//RootCAs: roots,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: "+err.Error())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(conn, in)
		wg.Done()
	}()

	go func() {
		io.Copy(out, conn)
		wg.Done()
	}()

	wg.Wait()

	conn.Close()
}
