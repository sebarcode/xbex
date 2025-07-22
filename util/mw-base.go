package util

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sebarcode/xbex/config"
)

func MwHttpAuth(ctx *kaos.Context, _ interface{}) (bool, error) {
	req := ctx.HttpRequest()
	if req == nil {
		return false, errors.New("missing http request")
	}
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return false, errors.New("missing authorization header")
	}
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return false, errors.New("invalid authorization header format")
	}
	tokenString := authHeader[len(bearerPrefix):]
	salt := config.Config().JwtSalt

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt), nil
	})
	if err != nil || !token.Valid {
		return false, errors.New("invalid token: " + err.Error())
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return false, errors.New("user_id not found in token")
	}
	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return false, errors.New("role not found in token")
	}

	if claims["exp"] == nil {
		return false, errors.New("token expiration not found")
	}
	if exp, ok := claims["exp"].(float64); !ok || exp < float64(time.Now().Unix()) {
		return false, errors.New("token expired")
	}

	ctx.Data().Set("jwt_token", tokenString)
	ctx.Data().Set("jwt_reference_id", userID)
	ctx.Data().Set("appuser_role", role)

	return true, nil
}

func MwCheckRole(roleid ...string) kaos.MWFunc {
	return func(ctx *kaos.Context, _ interface{}) (bool, error) {
		role := ctx.Data().Get("appuser_role", "").(string)
		if role == "" {
			return false, errors.New("missing role in context")
		}
		if !slices.Contains(roleid, role) {
			return false, fmt.Errorf("allowed role: %v", roleid)
		}
		return true, nil
	}
}

func MwLimitTake(limit int) kaos.MWFunc {
	return func(ctx *kaos.Context, payload interface{}) (bool, error) {
		if limit <= 0 {
			return false, errors.New("invalid limit value")
		}

		smPath := ctx.Data().Get("path", "").(string)
		if !(strings.HasSuffix(smPath, "/gets") || strings.HasSuffix(smPath, "/find")) {
			return true, nil
		}

		qp, ok := payload.(*dbflex.QueryParam)
		if !ok {
			return false, errors.New("payload is not a QueryParam")
		}
		take := qp.Take
		if take > limit || take == 0 {
			qp.Take = limit
		}
		val := reflect.ValueOf(payload)
		val.Elem().Set(reflect.ValueOf(qp).Elem())

		return true, nil
	}
}
