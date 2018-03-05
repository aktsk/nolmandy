package main

import (
	"flag"

	"github.com/aktsk/nolmandy/server"
)

func main() {
	const defaultPort = 8000
	var port int
	flag.IntVar(&port, "port", defaultPort, "Port to listen")
	flag.Parse()
	server.Serve(port)
}
