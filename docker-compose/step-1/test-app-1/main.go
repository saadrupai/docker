package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server2Port := os.Getenv("PORT2")
	if server2Port == "" {
		server2Port = "8090"
	}

	// Define the handler for the API
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from " + r.URL.Path)
		fmt.Fprintln(w, "Hello from Go server!")
	})

	http.HandleFunc("/hello-from-test2", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://test-app-2:" + server2Port + "/hello")
		if err != nil {
			http.Error(w, "Error calling test2 service", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response from test2 service", http.StatusInternalServerError)
			return
		}

		// Print the response body to the server logs
		fmt.Println("Response from test2 service:", string(body))
		w.Write(body)
	})

	// Start the server on port 8080
	fmt.Printf("Server running on port:%s\n", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		panic(err)
	}
}
