package logic

import (
	"errors"
	"time"

	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/datahub"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sebarcode/codekit"
	"github.com/sebarcode/xbex/config"
	"github.com/sebarcode/xbex/model"
)

type AuthHandler struct {
}

func (obj *AuthHandler) HttpAuth(ctx *kaos.Context, payload *string) (string, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return "", errors.New("missingDBConn")
	}

	req := ctx.HttpRequest()
	userName, pass, ok := req.BasicAuth()
	if !ok {
		return "", errors.New("missingAuth")
	}
	shaKey := config.Config().ShaKey
	saltPass := codekit.ShaString(pass, shaKey)

	meta, err := datahub.GetByID(h, new(model.AppUserMeta), userName)
	if err != nil {
		return "", errors.New("invalidCredential")
	}
	if meta.Password != saltPass {
		return "", errors.New("invalidCredential")
	}
	user, err := datahub.GetByID(h, new(model.AppUser), userName)
	if err != nil {
		return "", errors.New("invalidCredential")
	}
	if user.Status != "Active" {
		return "", errors.New("userNotActive")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     jwt.NewNumericDate(time.Now().Add(3 * time.Hour)), // token expires in 3 hours
	}

	signMethod := jwt.GetSigningMethod("HS256")
	token := jwt.NewWithClaims(signMethod, claims)
	jwtStr, err := token.SignedString([]byte(config.Config().JwtSalt))
	if err != nil {
		return "", errors.New("jwtSignFailed")
	}
	return jwtStr, nil
}
