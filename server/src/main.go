package main

import "net/http"

func main() {
	http.HandleFunc("/save", saveBatchData)
}


func saveBatchData(w http.ResponseWriter, r *http.Request) {

}