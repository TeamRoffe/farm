package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	// Apparently the way to do it
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/teamroffe/farm/pkg/drinks"
	"github.com/teamroffe/farm/pkg/pumps"
	"gopkg.in/ini.v1"
)

//FarmServer our main webserver package
type FarmServer struct {
	Status   *Status
	DB       *sql.DB
	Config   *ini.File
	PM       *pumps.PumpManager
	stopChan chan bool
}

//Status holds F.A.R.M status
type Status struct {
	mux      sync.Mutex
	Pouring  bool
	LastPour *time.Time
}

// farmResponse is the generic response message type
type farmResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//drinkInfo
type drinkInfo struct {
	Info        *drinks.Drink             `json:"info"`
	Ingredients *[]drinks.DrinkIngredient `json:"ingredients"`
}

//drinkResponse is the response format for /v1/drink/:id
type drinkResponse struct {
	ID          *int    `json:"id"`
	Name        *string `json:"drink_name"`
	Description *string `json:"description"`
	URL         *string `json:"url"`
	Ingredients []*drinks.DrinkIngredient
}

//handleHealthz dummy endpoint for health checking
func (server *FarmServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "Ok\n")
}

//pourOn toggles the global pouring flag making F.A.R.M DENY pour requests
func (server *FarmServer) pourON() {
	server.Status.mux.Lock()
	server.Status.Pouring = true
	server.Status.mux.Unlock()
}

//pourOff toggles the global pouring flag making F.A.R.M ACCEPT pour requests
func (server *FarmServer) pourOFF() {
	server.Status.mux.Lock()
	server.Status.Pouring = false
	now := time.Now()
	server.Status.LastPour = &now
	server.Status.mux.Unlock()
}

//Stop the F.A.R.M server
func (server *FarmServer) Stop() {
	defer close(server.stopChan)
	server.stopChan <- false
}

//Run starts the server
func (server *FarmServer) Run() error {
	db, err := sql.Open("mysql", server.getDSN())
	server.DB = db
	if err != nil {
		return err
	}
	defer db.Close()

	//Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	//Start pump manager
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.PM.Run()
	}()

	router := mux.NewRouter()

	//add our HTTP routes
	router.HandleFunc("/v1/categories", server.getCategories).Methods("GET", "POST")
	router.HandleFunc("/v1/drinks", server.handleDrink).Methods("GET")
	router.HandleFunc("/v1/drink/{id}", server.handleDrink).Methods("GET", "POST")
	router.HandleFunc("/v1/liquids", server.handleLiquid).Methods("GET")
	router.HandleFunc("/v1/liquid/{id}", server.handleLiquid).Methods("GET", "POST")
	router.HandleFunc("/v1/pour/{id}", server.handlePour).Methods("GET", "POST")
	router.HandleFunc("/healthz", server.handleHealthz).Methods("GET")
	router.HandleFunc("/v1/ports", server.handlePorts).Methods("GET")
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", server.Config.Section("server").Key("http_port").String()), router); err != nil {
			glog.Fatal(err)
		}
	}()

	<-server.stopChan
	server.PM.Stop()
	wg.Wait()
	return nil
}

// NewServer new farm pouring client
func NewServer() *FarmServer {
	var gracefulStop = make(chan bool)
	//Load our config file
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		glog.Fatalf("Fail to read file: %v", err)
		os.Exit(1)
	}

	pm, err := pumps.NewPumpManager(cfg)
	if err != nil {
		glog.Fatal(err)
	}

	now := time.Now()

	return &FarmServer{
		PM:     pm,
		Config: cfg,
		Status: &Status{
			LastPour: &now,
		},
		stopChan: gracefulStop,
	}
}
