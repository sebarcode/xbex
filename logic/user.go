package logic

import (
	"errors"

	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/datahub"
	"github.com/sebarcode/codekit"
	"github.com/sebarcode/xbex/model"
)

type UserHandler struct {
}

func (obj *UserHandler) ChangePassword(db *datahub.Hub, userid, password, salt string) error {
	userMeta, err := datahub.GetByID(db, new(model.AppUserMeta), userid)
	if err != nil {
		userMeta = &model.AppUserMeta{
			ID: userid,
		}
	}

	userMeta.Password = codekit.ShaString(password, salt)
	return db.Save(userMeta)
}

type CreateUserRequest struct {
	User     *model.AppUser
	Password string
}

func (obj *UserHandler) Create(ctx *kaos.Context, payload *CreateUserRequest) (*model.AppUser, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return nil, errors.New("missing: db")
	}

	if payload.User == nil {
		return nil, errors.New("missing: user info")
	}

	if payload.Password == "" {
		return nil, errors.New("missing: password")
	}

	salt := ctx.Data().Get("service_secret", "").(string)
	if salt == "" {
		return nil, errors.New("missing: service salt")
	}

	e := h.Insert(payload.User)
	if e != nil {
		return nil, e
	}

	if err := obj.ChangePassword(h, payload.User.ID, payload.Password, salt); err != nil {
		return nil, err
	}

	return payload.User, nil
}
