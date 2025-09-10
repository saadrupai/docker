package main

import (
	"fmt"
	"io"
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
		port = "8080"
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Test App 2\n"))
	})

	port2 := os.Getenv("PORT2")
	if port2 == "" {
		port2 = "8080"
	}

	app2Host := os.Getenv("APP2_HOST")
	if app2Host == "" {
		app2Host = "localhost"
	}

	app2URL := makeApp3Url(app2Host, port2)
	fmt.Println(app2URL)

	http.HandleFunc("/hello-from-app3", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(app2URL)
		if err != nil {
			http.Error(w, "Failed to reach App 2", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "App 2 returned an error", http.StatusInternalServerError)
			return
		}

		respByte, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response from App 2", http.StatusInternalServerError)
			return
		}

		fullResponse := fmt.Sprintf("Hello from Test App 1. Also, App 2 says: %s", string(respByte))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fullResponse))
	})

	fmt.Println("Server starting at port: ", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("failed to start server at port: ", port)
		panic(err)
	}
}

func makeApp3Url(app2Host string, port2 string) string {
	app2URL := fmt.Sprintf("http://%s:%s/hello", app2Host, port2)
	return app2URL
}
