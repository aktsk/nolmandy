package main

import (
	"net/http"

	"github.com/aktsk/nolmandy/server"
)

func init() {
	http.HandleFunc("/", server.Handler)
}
