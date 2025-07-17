package api

import (
	"github.com/ariefdarmawan/suim"
	"github.com/sebarcode/dbmod"
	"github.com/sebarcode/rayiapp"
	"github.com/sebarcode/xbex/logic"
	"github.com/sebarcode/xbex/model"
	"github.com/sebarcode/xbex/util"
)

func RegisterCore(app *rayiapp.App) error {
	modDB := dbmod.New()
	modUI := suim.New()

	s := app.Service()
	s.Group().SetMod(modUI).
		AllowOnlyRoute("gridconfig", "formconfig").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user"),
			s.RegisterModel(&model.AppUserCreate{}, "user-create").AllowOnlyRoute("formconfig"),
		)

	s.Group().SetMod(modDB).
		RegisterMWs(util.MwHttpAuth).
		AllowOnlyRoute("get", "gets", "find").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user"),
		)

	s.Group().SetMod(modUI).
		RegisterMWs(util.MwHttpAuth).
		AllowOnlyRoute("insert", "update", "delete", "save").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user").AllowOnlyRoute("update").RegisterMWs(util.MwCheckRole("Admin")),
		)

	s.RegisterModel(new(logic.AuthHandler), "rbac")

	s.Data().Set("jwt_salt", app.Config.Data.GetString("jwt_salt"))

	return nil
}
