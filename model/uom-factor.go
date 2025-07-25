package model

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitFactor struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id" key:"1" form_read_only_edit:"1" form_section:"General" form_section_auto_col:"2"`
	FromUnitID        string
	ToUnitID          string
	FromQty           float64
	ToQty             float64
	Created           time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info" form_section_auto_col:"2"`
	LastUpdate        time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info"`
}

func (o *UnitFactor) TableName() string {
	return "UnitFactors"
}

func (o *UnitFactor) FK() []*orm.FKConfig {
	return orm.DefaultRelationManager().FKs(o)
}

func (o *UnitFactor) ReverseFK() []*orm.ReverseFKConfig {
	return orm.DefaultRelationManager().ReverseFKs(o)
}

func (o *UnitFactor) SetID(keys ...interface{}) {
	o.ID = keys[0].(string)
}

func (o *UnitFactor) GetID(dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{o.ID}
}

func (o *UnitFactor) PreSave(dbflex.IConnection) error {
	if o.ID == "" {
		o.ID = primitive.NewObjectID().Hex()
	}
	if o.Created.IsZero() {
		o.Created = time.Now()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *UnitFactor) PostSave(dbflex.IConnection) error {
	return nil
}
func (o *UnitFactor) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{
		{Name: "idx_from_unit_id", Fields: []string{"FromUnitID", "ToUnitID"}},
	}
}
