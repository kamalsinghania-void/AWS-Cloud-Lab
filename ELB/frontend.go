package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
)

type PageData struct {
	Message string
}

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Send a connection request to the load balancer
		conn, err := net.Dial("tcp", "<Provide NLB DNS name here>:6000")
		if err != nil {
			log.Println("We failed to connect to the back-end server")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Set the value we'll send to the backend
		dynamodbValue := "client-01"
		_, err = conn.Write([]byte(dynamodbValue))
		if err != nil {
			log.Println("Error writing to connection:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Receive data from backend
		recvdSlice := make([]byte, 150)
		n, err := conn.Read(recvdSlice)
		if err != nil {
			log.Println("Error reading from connection:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data := PageData{
			Message: string(recvdSlice[:n]),
		}

		// Create an HTML template
		tmpl, err := template.New("index").Parse(`
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>Hello, Go!</title>
                <style>
                    body {
                        font-family: fantasy;
                        background-color: #17202D;
                        margin: 0;
                        display: flex;
                        justify-content: center;
                        align-items: center;
                        height: 100vh;
                    }
                    .container {
                        text-align: center;
                    }
                    .message {
                        font-weight: bold;
                        color: lightyellow;
                    }
                </style>
            </head>
            <body>
                <div class="container">
                    <h1 class="message">Congratulations!</h1>
                    <h2 class="message">You've connected to the back-end server and fetched the following data</h2>
                    <h3 class="message">{{.Message}}</h3>
                </div>
            </body>
            </html>
        `)
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the template with the provided data
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	// Register the handler function for the root URL path "/"
	http.HandleFunc("/", handler)

	// Start the HTTP server on port 8080
	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}