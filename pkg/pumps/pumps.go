package pumps

import (
	"fmt"
	"os"
	"time"

	"github.com/cloudflare/cfssl/log"
	rpio "github.com/stianeikeland/go-rpio"
)

//PumpMSG holds the pour job
type PumpMSG struct {
	Port int
	Time time.Duration
}

//PumpManager handles our pour jobs
type PumpManager struct {
	Queue     chan *PumpMSG
	PumpPorts []*pumpPort
	RpiHW     bool
	quitChan  chan bool
}

type pumpPort struct {
	Pin      rpio.Pin
	LiquidID int
}

// NewPumpManager returns new PumpManager
func NewPumpManager() (*PumpManager, error) {
	rpi := false
	if _, err := os.Stat("/dev/mem"); !os.IsNotExist(err) {
		rpi = true
	}
	msg := make(chan *PumpMSG)
	return &PumpManager{
		Queue: msg,
		RpiHW: rpi,
	}, nil
}

//Run starts the pumpmanager
func (pm *PumpManager) Run() error {
	pm.quitChan = make(chan bool)
	if pm.RpiHW {
		defer rpio.Close()

		if err := rpio.Open(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		Relay1 := rpio.Pin(18)
		Relay1.Output()

		Relay2 := rpio.Pin(22)
		Relay2.Output()

		pm.PumpPorts = append(pm.PumpPorts, &pumpPort{Pin: Relay1, LiquidID: 1}, &pumpPort{Pin: Relay2, LiquidID: 2})
	}
	jobDone := make(chan int)
	defer close(jobDone)
	for {
		select {
		case message := <-pm.Queue:
			go pm.pour(message, jobDone)
		case ch := <-jobDone:
			log.Infof("Done: %d", ch)
		case <-pm.quitChan:
			return nil
		}
	}
}

func (pm *PumpManager) Stop() error {
	close(pm.quitChan)
	return nil
}

func (pm *PumpManager) pour(job *PumpMSG, done chan<- int) error {
	log.Infof("pour: %s", job)
	time.Sleep(job.Time)
	done <- job.Port
	return nil
}
