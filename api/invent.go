package api

import (
	"errors"

	"github.com/ariefdarmawan/suim"
	"github.com/sebarcode/dbmod"
	"github.com/sebarcode/rayiapp"
	"github.com/sebarcode/xbex/logic"
	"github.com/sebarcode/xbex/model"
	"github.com/sebarcode/xbex/util"
)

func RegisterInvent(app *rayiapp.App) error {
	modUI := suim.New()

	s := app.Service()
	jwtSalt := app.Config.Data.GetString("jwt_salt")
	if jwtSalt == "" {
		return errors.New("missing jwt_salt in config")
	}
	s.Data().Set("jwt_salt", jwtSalt)
	s.RegisterMW(util.MwLimitTake(100), "invent_limit_take")

	s.Group().SetMod(modUI).
		AllowOnlyRoute("gridconfig", "formconfig").
		Apply(
			s.RegisterModel(&model.Item{}, "item"),
			s.RegisterModel(&model.ItemStock{}, "stock"),
			s.RegisterModel(&model.Unit{}, "unit"),
			s.RegisterModel(&model.UnitFactor{}, "unit-factor"),
			s.RegisterModel(&model.JournalHeader{}, "journal-header"),
			s.RegisterModel(&model.JournalLine{}, "journal-line"),
		)

	modDB := dbmod.New()
	s.Group().SetMod(modDB).
		RegisterMWs(util.MwHttpAuth).
		AllowOnlyRoute("get", "gets", "find").
		Apply(
			s.RegisterModel(&model.Item{}, "item"),
			s.RegisterModel(&model.ItemStock{}, "stock"),
			s.RegisterModel(&model.Unit{}, "unit"),
			s.RegisterModel(&model.UnitFactor{}, "unit-factor"),
			s.RegisterModel(&model.JournalHeader{}, "journal-header"),
			s.RegisterModel(&model.JournalLine{}, "journal-line"),
		)

	s.Group().SetMod(modDB).
		RegisterMWs(
			util.MwHttpAuth,
			util.MwCheckRole("Admin", "Manager")).
		AllowOnlyRoute("insert", "update", "delete", "save").
		Apply(
			s.RegisterModel(&model.Item{}, "item"),
			s.RegisterModel(&model.ItemStock{}, "stock"),
			s.RegisterModel(&model.Unit{}, "unit"),
			s.RegisterModel(&model.UnitFactor{}, "unit-factor"),
		)

	s.Group().SetMod(modDB).
		RegisterMWs(util.MwHttpAuth, util.MwCheckRole()).
		AllowOnlyRoute("insert", "update", "delete", "save").
		Apply(
			s.RegisterModel(&model.JournalHeader{}, "journal-header"),
			s.RegisterModel(&model.JournalLine{}, "journal-line"),
		)

	s.Group().
		RegisterMWs(util.MwHttpAuth, util.MwCheckRole()).
		Apply(
			s.RegisterModel(new(logic.JournalHandler), "journal").AllowOnlyRoute("submit"),
		)

	s.Group().
		RegisterMWs(util.MwHttpAuth, util.MwCheckRole("Admin", "Manager")).
		Apply(
			s.RegisterModel(new(logic.JournalHandler), "journal").DisableRoute("submit"),
		)

	s.RegisterModel(new(model.JournalLog), "journal/log").SetMod(modUI, modDB).AllowOnlyRoute("gets", "formconfig", "gridconfig")
	return nil
}
