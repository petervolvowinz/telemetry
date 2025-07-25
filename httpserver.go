package main

import (
	"log"
	"net/http"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fileserver)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
