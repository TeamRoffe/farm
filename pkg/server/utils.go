package server

import "fmt"

func (server *FarmServer) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		server.Config.Section("database").Key("username").String(),
		server.Config.Section("database").Key("password").String(),
		server.Config.Section("database").Key("hostname").String(),
		server.Config.Section("database").Key("port").String(),
		server.Config.Section("database").Key("database").String(),
	)
}
