package api

import (
	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/suim"
	"github.com/sebarcode/dbmod"
	"github.com/sebarcode/xbex/model"
	"github.com/sebarcode/xbex/util"
)

func RegisterCore(s *kaos.Service) error {
	modDB := dbmod.New()
	modUI := suim.New()

	s.Group().SetMod(modUI).
		AllowOnlyRoute("gridconfig", "formconfig").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user"),
			s.RegisterModel(&model.AppUserCreate{}, "user-create").AllowOnlyRoute("formconfig"),
		)

	s.Group().SetMod(modDB).
		RegisterMWs(util.MwAuth).
		AllowOnlyRoute("get", "gets", "find").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user"),
		)

	s.Group().SetMod(modUI).
		RegisterMWs(util.MwAuth).
		AllowOnlyRoute("insert", "update", "delete", "save").
		Apply(
			s.RegisterModel(&model.AppUser{}, "user").AllowOnlyRoute("update").RegisterMWs(util.MwCheckRole("Admin")),
		)

	return nil
}
