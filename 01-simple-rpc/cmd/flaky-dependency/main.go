package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	handler, err := FromEnv()
	if err != nil {
		fmt.Println("failed to load handler from env")
		os.Exit(2)
	}
	http.HandleFunc("/do-work", handler.Handle)
	http.HandleFunc("/set-flakiness", handler.SetFlakiness)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		os.Exit(3)
	}
}
