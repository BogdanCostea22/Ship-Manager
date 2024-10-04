package main

import (
	"Ship_Manager/internal/server"
	"fmt"
	"log"
)

func main() {
	server := server.NewServer()

	log.Printf("Server is running at address %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
