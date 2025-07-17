package util

import "git.kanosolution.net/kano/kaos"

func MwAuth(ctx *kaos.Context, _ interface{}) (bool, error) {
	return true, nil
}

func MwCheckRole(roleid string) kaos.MWFunc {
	return func(ctx *kaos.Context, _ interface{}) (bool, error) {
		return true, nil
	}
}
