package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/teamroffe/farm/pkg/server"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: farm -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Parse()
}

// our main function
func main() {
	glog.Info("Starting F.A.R.M")

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	server := server.NewServer()

	go func() {
		sig := <-gracefulStop
		glog.Infof("Caught sig: %+v Wait to finish processing", sig)
		server.Stop()
	}()

	err := server.Run()
	if err != nil {
		glog.Fatal(err)
	}
	os.Exit(0)
}
