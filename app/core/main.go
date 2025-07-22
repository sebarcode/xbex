package main

import (
	"flag"
	"fmt"

	"git.kanosolution.net/kano/kaos"
	_ "github.com/ariefdarmawan/flexpg"
	_ "github.com/kanoteknologi/hd"
	"github.com/sebarcode/codekit"
	_ "github.com/sebarcode/htev"
	"github.com/sebarcode/rayiapp"
	"github.com/sebarcode/xbex/api"
)

var (
	appName         = "core-rest"
	basePoint       = "api/core"
	version         = "v1.0.0"
	serviceNameRest = "core-rest"
	serviceNameEv   = "core-event"
	app             *rayiapp.App
	logWriter       = kaos.CreateLogWithPrefix("core-rest")
	configFile      = flag.String("config", "app.yml", "path to config file")
)

func main() {
	var err error

	flag.Parse()
	codekit.DefaultCase = codekit.CaseAsIs
	kaos.NamingType = kaos.NamingIsLower

	logWriter.Infof("initiating %s %s", appName, version)
	app, err = rayiapp.CreateApp(*configFile, &rayiapp.AppOpts{
		Logger:     logWriter,
		Apps:       []string{serviceNameRest, serviceNameEv},
		Publishers: []string{serviceNameEv},
	})

	app.Service().SetBasePoint(basePoint)
	app.Name = fmt.Sprintf("%s %s", appName, version)
	if err != nil {
		app.Exit(-1, err.Error())
	}

	if e := api.RegisterCore(app); e != nil {
		app.Exit(-1, e.Error())
	}

	if err := app.Start(); err != nil {
		app.Exit(-1, err.Error())
	}
	app.WaitForGraceShutdown()
}
