package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Test App 3\n"))
	})

	fmt.Println("Server starting at port: ", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("failed to start server at port: ", port)
		panic(err)
	}
}
