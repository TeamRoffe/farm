package main

import (
	"github.com/teamroffe/farm/pkg/server"
)

// our main function
func main() {
	server := server.NewServer()
	server.Run()
}
