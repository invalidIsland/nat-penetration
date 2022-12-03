package main

import (
	"encoding/json"
	"log"
	"nat-penetration/define"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		bytes, err := json.Marshal(query)
		if err != nil {
			log.Printf("Marshal Error: %v", err)
		}
		writer.Write(bytes)
	})
	log.Printf("LocalServer has started: %s\n", define.LocalServerAddr)
	http.ListenAndServe(define.LocalServerAddr, nil)
}
