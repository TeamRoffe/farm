package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"

	// Apparently the way to do it
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/teamroffe/farm/pkg/drinks"
	"github.com/teamroffe/farm/pkg/pumps"
	"gopkg.in/ini.v1"
)

// FarmServer our main webserver package
type FarmServer struct {
	Status Status
	DB     *sql.DB
	Config *ini.File
	PM     *pumps.PumpManager
}

// Status holds F.A.R.M status
type Status struct {
	mux     sync.Mutex
	Pouring bool
}

type farmResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (server *FarmServer) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "Ok\n")
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

	router.HandleFunc("/v1/categories", server.getCategories).Methods("GET", "POST")
	router.HandleFunc("/v1/drinks", server.handleDrink).Methods("GET")
	router.HandleFunc("/v1/drink/{id}", server.handleDrink).Methods("GET", "POST")
	router.HandleFunc("/v1/liquids", server.handleLiquid).Methods("GET")
	router.HandleFunc("/v1/liquid/{id}", server.handleLiquid).Methods("GET", "POST")
	router.HandleFunc("/v1/pour/{id}", server.pour).Methods("GET", "POST")
	router.HandleFunc("/healthz", server.healthz).Methods("GET")
	router.HandleFunc("/v1/status/pump/{pump}", server.pumpStatus).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%s", server.Config.Section("server").Key("http_port").String()), router)
}

// NewServer new farm pouring client
func NewServer() *FarmServer {

	cfg, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	pm, err := pumps.NewPumpManager()
	if err != nil {
		panic(err)
	}
	go pm.Run()

	return &FarmServer{
		PM:     pm,
		Config: cfg,
	}
}
