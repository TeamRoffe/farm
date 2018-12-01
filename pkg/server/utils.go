package server

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/teamroffe/farm/pkg/drinks"
)

func (server *FarmServer) hasLiquids(ingredients []*drinks.DrinkIngredient) error {
	var found bool
	for _, ingr := range ingredients {
		found = false
		for _, port := range server.PM.Ports {
			if *port.LiquidID == *ingr.LiquidID {
				glog.Infof("We got liquid %d localy", *ingr.LiquidID)
				found = true
				continue
			}
		}
		if !found {
			return fmt.Errorf("This F.A.R.M does not carry liquid %d", *ingr.LiquidID)
		}
	}
	return nil
}

func (server *FarmServer) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		server.Config.Section("database").Key("username").String(),
		server.Config.Section("database").Key("password").String(),
		server.Config.Section("database").Key("hostname").String(),
		server.Config.Section("database").Key("port").String(),
		server.Config.Section("database").Key("database").String(),
	)
}
