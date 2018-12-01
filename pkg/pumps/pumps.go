package pumps

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	rpio "github.com/stianeikeland/go-rpio"
	"gopkg.in/ini.v1"
)

//PumpMSG holds the pour job
type PumpMSG struct {
	LiquidID *int
	Time     time.Duration
	Done     chan bool
}

//PumpManager handles our pour jobs
type PumpManager struct {
	Cfg      *ini.File
	Ports    []*pumpPort
	Queue    chan *PumpMSG
	RpiHW    bool
	quitChan chan bool
}

//pumpPort holds what pin on the pi a specific beverage sits
type pumpPort struct {
	Pin      *rpio.Pin `json:"pin"`
	LiquidID *int      `json:"liquid_id"`
}

// NewPumpManager returns new PumpManager
func NewPumpManager(cfg *ini.File) (*PumpManager, error) {
	var ports []*pumpPort

	var RpiHW bool
	if _, err := os.Stat("/dev/mem"); !os.IsNotExist(err) {
		RpiHW = true
	}

	msg := make(chan *PumpMSG)

	return &PumpManager{
		Cfg:   cfg,
		Ports: ports,
		RpiHW: RpiHW,
		Queue: msg,
	}, nil
}

//Run starts the pumpmanager
func (pm *PumpManager) Run() error {
	pm.quitChan = make(chan bool)
	jobDone := make(chan *int)
	defer close(jobDone)

	if pm.RpiHW {
		glog.Info("Running on rpi hardware")
		defer rpio.Close()

		if err := rpio.Open(); err != nil {
			glog.Fatal(err)
			os.Exit(1)
		}

		amount, err := pm.Cfg.Section("ports").Key("no_ports").Int()
		if err != nil {
			return err
		}

		for i := 1; i < (amount + 1); i++ {
			port := pm.Cfg.Section("ports").Key(fmt.Sprintf("%d", i)).String()
			glog.Infof("Portconfig: %s", port)
			parts := strings.Split(port, ":")
			pinID, err := strconv.Atoi(parts[0])
			if err != nil {
				return err
			}
			liquidID, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}
			pin := rpio.Pin(pinID)
			pin.Output()
			pumpPort := &pumpPort{
				Pin:      &pin,
				LiquidID: &liquidID,
			}

			pm.Ports = append(pm.Ports, pumpPort)

		}
	}
	defer glog.Info("PM exited")
	for {
		select {
		case message := <-pm.Queue:
			go pm.pour(message, jobDone)
		case ch := <-jobDone:
			glog.Infof("Done liquid: %d", *ch)
		case <-pm.quitChan:
			glog.Info("Stopping pump manager")
			return nil
		}
	}
}

//Stop the pumpmanager
func (pm *PumpManager) Stop() {
	pm.quitChan <- false
	close(pm.quitChan)
}

//pour handles the port controll of the rpi
func (pm *PumpManager) pour(job *PumpMSG, done chan<- *int) error {
	defer close(job.Done)
	glog.Infof("PM Pour: %d %s", *job.LiquidID, job.Time)
	for _, port := range pm.Ports {
		if *port.LiquidID == *job.LiquidID {
			port.Pin.High()
			time.Sleep(job.Time)
			port.Pin.Low()
		}
	}
	done <- job.LiquidID
	return nil
}
