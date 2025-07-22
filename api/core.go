package api

import (
	"errors"

	"github.com/ariefdarmawan/serde"
	"github.com/ariefdarmawan/suim"
	"github.com/kanoteknologi/hd"
	"github.com/sebarcode/dbmod"
	"github.com/sebarcode/rayiapp"
	"github.com/sebarcode/xbex/config"
	"github.com/sebarcode/xbex/logic"
	"github.com/sebarcode/xbex/model"
	"github.com/sebarcode/xbex/util"
)

func RegisterCore(app *rayiapp.App) error {
	modConfig := new(config.ModConfig)
	err := serde.Serde(app.Config.Data, modConfig)
	if err != nil {
		return err
	}
	if modConfig.JwtSalt == "" {
		return errors.New("missing jwt_salt in config")
	}
	if modConfig.ShaKey == "" {
		return errors.New("missing sha_secret in config")
	}
	config.SetConfig(modConfig)

	modDB := dbmod.New()
	modUI := suim.New()

	s := app.Service()
	s.Group().
		SetDeployer(hd.DeployerName).
		SetMod(modUI).
		AllowOnlyRoute("gridconfig", "formconfig").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user"),
			s.RegisterModel(&model.AppUserCreate{}, "user-create").AllowOnlyRoute("formconfig"),
		)

	s.Group().
		SetDeployer(hd.DeployerName).
		SetMod(modDB).
		RegisterMWs(util.MwHttpAuth).
		AllowOnlyRoute("get", "gets", "find").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user"),
		)

	s.Group().SetMod(modUI).
		SetDeployer(hd.DeployerName).
		RegisterMWs(util.MwHttpAuth).
		AllowOnlyRoute("insert", "update", "delete", "save").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user").AllowOnlyRoute("update").RegisterMWs(util.MwCheckRole("Admin")),
		)

	s.RegisterModel(new(logic.AuthHandler), "rbac").SetDeployer(hd.DeployerName)

	s.Data().Set("jwt_salt", app.Config.Data.GetString("jwt_salt"))

	return nil
}
