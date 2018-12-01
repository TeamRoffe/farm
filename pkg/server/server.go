package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	rpio "github.com/stianeikeland/go-rpio"
	"github.com/teamroffe/farm/pkg/drinks"
	"gopkg.in/ini.v1"
)

// FarmServer our main webserver package
type FarmServer struct {
	ListenPort uint16
	Status     Status
	DB         *sql.DB
	Config     *ini.File
	Relay1     rpio.Pin
}

// Status holds F.A.R.M status
type Status struct {
	mux     sync.Mutex
	Pouring bool
}

type pourResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (server *FarmServer) healthz(w http.ResponseWriter, r *http.Request) {
	return
}

func (server *FarmServer) getingredients(drinkID int) []drinks.DrinkIngredient {
	var ingredients []drinks.DrinkIngredient
	results, err := server.DB.Query("SELECT id, drink_id, liquid_id, volume FROM drink_ingredients WHERE drink_id = ?", drinkID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()
	for results.Next() {
		var ingredient drinks.DrinkIngredient
		// for each row, scan the result into our tag composite object
		err = results.Scan(&ingredient.ID, &ingredient.DrinkID, &ingredient.LiquidID, &ingredient.Volume)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		} else {
			ingredients = append(ingredients, ingredient)
		}
		// and then print out the tag's Name attribute
	}
	return ingredients
}

func (server *FarmServer) pumpStatus(w http.ResponseWriter, r *http.Request) {
	return
}

type drinkInfo struct {
	Info        *drinks.Drink             `json:"info"`
	Ingredients *[]drinks.DrinkIngredient `json:"ingredients"`
}

//Run starts the server
func (server *FarmServer) Run() error {
	defer rpio.Close()
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server.Relay1 = rpio.Pin(18)
	server.Relay1.Output()

	db, err := sql.Open("mysql", server.getDSN())
	server.DB = db
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	router := mux.NewRouter()
	router.HandleFunc("/v1/drinkid/{id}", server.getDrink).Methods("GET")
	router.HandleFunc("/healthz", server.healthz).Methods("GET")
	router.HandleFunc("/v1/pour/{id}", server.pour).Methods("GET")
	router.HandleFunc("/v1/status/pump/{pump}", server.pumpStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.ListenPort), router))
	return nil
}

// NewServer new farm pouring client
func NewServer() *FarmServer {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return &FarmServer{
		ListenPort: 8000,
		Config:     cfg,
	}
}
