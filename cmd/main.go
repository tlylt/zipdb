package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tlylt/zipdb/domain"
	"github.com/tlylt/zipdb/zipdb"
)

type locationByZip struct {
	zipDB map[string]domain.Location
}

func (l locationByZip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zip := r.URL.Path[1:]
	location, ok := l.zipDB[zip]
	if !ok {
		http.Error(w, "Zip not found", http.StatusNotFound)
		return
	}
	b, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Location: %v", string(b))
}

func main() {
	zipDB := zipdb.LoadZipDB()

	fmt.Println("Starting the web server at http://localhost:8080")
	err := http.ListenAndServe(":8080", locationByZip{zipDB: zipDB})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting...")
	}
}
