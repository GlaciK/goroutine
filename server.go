package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"main/tools"
	"net/http"
)

func middleware(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ngecek apakah method yang direquest itu sesuai dengan di mux
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// kalau sesuai nanti dia langsung ngarahuin ke function handler nya
		handlerFunc(w, r)
	}
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hellooo !")
}

func sendFileHandler(w http.ResponseWriter, r *http.Request) {
	jsonData := r.FormValue("Person")

	var person data.Person
	err := json.Unmarshal([]byte(jsonData), &person)
	tools.ErrorHandler(err)

	fmt.Println("JSON : ", person)

	// handler berisi nama file, size file, dsb nya
	// file : yang bakal berisi semua data
	file, handler, err := r.FormFile("File")
	tools.ErrorHandler(err)

	fmt.Printf("Recived file: %s\n", handler.Filename)

	fileContent, err := io.ReadAll(file)
	tools.ErrorHandler(err)

	fmt.Printf("File Content: \n%s\n", fileContent)

	fmt.Fprintln(w, "Successfully recived data")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware(http.MethodGet, getMessageHandler))
	mux.HandleFunc("/sendFile", middleware(http.MethodPost, sendFileHandler))

	server := http.Server{
		Addr:    "localhost:9876",
		Handler: mux,
	}

	err := server.ListenAndServe()
	tools.ErrorHandler(err)
}
