package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Healthcheck(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		getHealthcheck(res, req)
	} else {
		http.NotFound(res, req)
	}
}

var count int = 0

func getHealthcheck(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	io.WriteString(res, fmt.Sprint(count))
	count++
}

func main() {
	var port string = os.Getenv("PORT")
	if len(port) == 0 {
		port = ":3000"
	}

	fmt.Println("Listening on port:" + port)

	http.Handle("/", http.HandlerFunc(Healthcheck))

	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
}
