package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/tlylt/zipdb/domain"
	"github.com/tlylt/zipdb/zipdb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type LocationDB struct {
	gorm.Model
	domain.Location `gorm:"embedded"`
}

func (l LocationDB) String() string {
	return fmt.Sprintf("ID: %d, Country: %s, Zip: %s, City: %s, StateLong: %s, State: %s, County: %s, Lat: %f, Long: %f",
		l.ID, l.Country, l.Zip, l.City, l.StateLong, l.State, l.County, l.Lat, l.Long)
}

type locationByZip struct {
	zipDB *gorm.DB
}

func (l locationByZip) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zip := r.URL.Path[1:]
	if zip == "favicon.ico" {
		return
	}

	var location LocationDB
	l.zipDB.First(&location, "Zip = ?", zip)

	fmt.Println(location)
	b, err := json.MarshalIndent(location, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Location: %v", string(b))
}

func main() {
	// Remove the existing database file to reset the database
	os.Remove("zipdb.db")

	db, err := gorm.Open(sqlite.Open("zipdb.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&LocationDB{})

	locations := zipdb.ParseLocations("US.txt")

	var locationDBs []LocationDB
	for _, location := range locations {
		locationDBs = append(locationDBs, LocationDB{Location: location})
	}

	// Load the data into the database
	db.CreateInBatches(&locationDBs, 100)

	// Start the web server
	fmt.Println("Starting the web server at http://localhost:8080")
	err = http.ListenAndServe(":8080", locationByZip{zipDB: db})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Exiting...")
	}
}
