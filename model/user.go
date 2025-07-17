package model

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppUserCreate struct {
	ID    string
	Name  string
	Email string
	Role  string `form_items:"User|Manager|Admin"`
}

type AppUser struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id" key:"1" form_read_only_edit:"1" form_section:"General" form_section_auto_col:"2"`
	Name              string
	Email             string
	Role              string    `form_items:"User|Manager|Admin"`
	Status            string    `form_items:"Draft|Active|Inactive"`
	Created           time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info" form_section_auto_col:"2"`
	LastUpdate        time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info"`
}

func (o *AppUser) TableName() string {
	return "AppUsers"
}

func (o *AppUser) FK() []*orm.FKConfig {
	return orm.DefaultRelationManager().FKs(o)
}

func (o *AppUser) ReverseFK() []*orm.ReverseFKConfig {
	return orm.DefaultRelationManager().ReverseFKs(o)
}

func (o *AppUser) SetID(keys ...interface{}) {
	o.ID = keys[0].(string)
}

func (o *AppUser) GetID(dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{o.ID}
}

func (o *AppUser) PreSave(dbflex.IConnection) error {
	if o.ID == "" {
		o.ID = primitive.NewObjectID().Hex()
	}
	if o.Created.IsZero() {
		o.Created = time.Now()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *AppUser) PostSave(dbflex.IConnection) error {
	return nil
}
func (o *AppUser) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{}
}
