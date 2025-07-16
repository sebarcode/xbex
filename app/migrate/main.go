package main

import (
	"git.kanosolution.net/kano/dbflex/orm"
	"github.com/ariefdarmawan/datahub"
	"github.com/sebarcode/logger"
	"github.com/sebarcode/xbex/model"
)

var (
	db *datahub.Hub
	lw *logger.LogEngine
)

func main() {
	models := []orm.DataModel{
		new(model.JournalHeader),
		new(model.JournalLine),
	}

	for _, model := range models {
		if err := db.EnsureDb(model); err != nil {
			panic(err)
		}
		lw.Infof("model %s has been migrated successfully", model.TableName())
	}
}
