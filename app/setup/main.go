package main

import (
	"flag"
	"os"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/datahub"
	_ "github.com/ariefdarmawan/flexpg"
	"github.com/sebarcode/codekit"
	"github.com/sebarcode/rayiapp"
	"github.com/sebarcode/xbex/logic"
	"github.com/sebarcode/xbex/model"
)

var (
	serviceName    = "setup"
	configFile     = flag.String("config", "app.yml", "path to config file")
	migrate        = flag.Bool("migrate", false, "run migration")
	changeRootPass = flag.String("pass", "", "new root password")
	logger         = kaos.CreateLogWithPrefix(serviceName)
)

func main() {
	flag.Parse()
	codekit.DefaultCase = codekit.CaseAsIs
	kaos.NamingType = kaos.NamingIsLower

	app, err := rayiapp.CreateApp(*configFile, &rayiapp.AppOpts{
		Logger:     logger,
		Apps:       []string{"x"},
		Publishers: []string{"x"},
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(-1)
	}

	err = app.StartDataHub()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(-1)
	}

	if *migrate {
		runMigration(app)
	}

	if *changeRootPass != "" {
		doChangeRootPass(app)
	}
}

func runMigration(app *rayiapp.App) func() {
	db, err := app.Service().GetDataHub("default", "")
	if err != nil {
		logger.Error("no db:" + err.Error())
		os.Exit(-1)
	}
	defer db.Close()

	conn, err := db.GetClassicConnection()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(-1)
	}
	defer conn.Close()
	dbflex.SetLogger(logger)

	models := []orm.DataModel{
		new(model.AppUser),
		new(model.AppUserMeta),
		new(model.Unit),
		new(model.Item),
		new(model.UnitFactor),
		new(model.JournalHeader),
		new(model.JournalLine),
	}

	for _, model := range models {
		logger.Infof("processing %s", model.TableName())
		if err := orm.EnsureDb(conn, model); err != nil {
			logger.Errorf("fail: %s, %s", model.TableName(), err.Error())
			continue
		}
		logger.Infof("done: %s", model.TableName())
	}

	return nil
}

func doChangeRootPass(app *rayiapp.App) error {
	db, err := app.Service().GetDataHub("default", "")
	if err != nil {
		logger.Error("no db:" + err.Error())
		os.Exit(-1)
	}
	defer db.Close()

	// root user
	if *changeRootPass != "" {
		sData := app.Config.Data
		salt := sData.GetString("secret")
		if salt == "" {
			logger.Error("missing: salt")
			return nil
		}

		logger.Infof("changing root password")
		user, err := datahub.GetByID(db, new(model.AppUser), "root")
		if err != nil {
			user = new(model.AppUser)
			user.ID = "root"
			user.Name = "Root User"
			user.Email = "root@app.com"
			user.Role = "Admin"
			user.Status = "Active"
			if err = db.Insert(user); err != nil {
				logger.Errorf("failed to create root user: %s", err.Error())
				return nil
			}
		}

		if err = new(logic.UserHandler).ChangePassword(db, user.ID, *changeRootPass, salt); err != nil {
			logger.Errorf("failed to change root password: %s", err.Error())
			return nil
		}
	}

	return nil
}
