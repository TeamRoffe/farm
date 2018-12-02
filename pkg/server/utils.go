package server

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/teamroffe/farm/pkg/drinks"
)

//createIng returns a pointer to the int
func createInt(x int) *int {
	return &x
}

//createString returns a pointer to the string
func createString(x string) *string {
	return &x
}

//farmResp returns a generic message objects from status & message
func farmResp(status int, message string) *farmResponse {
	return &farmResponse{
		Status:  createInt(status),
		Message: createString(message),
	}
}

//hasLiquids checks if the F.A.R.M has the ingredients to make the drink
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
			return fmt.Errorf("This F.A.R.M unit does not carry liquid %d", *ingr.LiquidID)
		}
	}
	return nil
}

//getDSN builds a DB DSN from our .ini config
func (server *FarmServer) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		server.Config.Section("database").Key("username").String(),
		server.Config.Section("database").Key("password").String(),
		server.Config.Section("database").Key("hostname").String(),
		server.Config.Section("database").Key("port").String(),
		server.Config.Section("database").Key("database").String(),
	)
}
