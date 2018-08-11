package main

import (
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/bandwidth/upload", saveData)
	http.ListenAndServe(":8080", nil)
}


func saveData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Data Received!")
}