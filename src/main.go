package main

import (
	"github.com/teamroffe/farm/pkg/server"
)

// our main function
func main() {
	server := server.NewServer()
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
