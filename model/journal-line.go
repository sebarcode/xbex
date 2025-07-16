package model

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JournalLine struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id" key:"1" form_read_only_edit:"1" form_section:"General" form_section_auto_col:"2"`
	JournalID         string
	LineNo            int
	LineRefID         string
	Name              string
	Qty               float64
	AmountEach        float64
	Amount            float64
	Description       string
	Created           time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info" form_section_auto_col:"2"`
	LastUpdate        time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info"`
}

func (o *JournalLine) TableName() string {
	return "JournalLines"
}

func (o *JournalLine) FK() []*orm.FKConfig {
	return orm.DefaultRelationManager().FKs(o)
}

func (o *JournalLine) ReverseFK() []*orm.ReverseFKConfig {
	return orm.DefaultRelationManager().ReverseFKs(o)
}

func (o *JournalLine) SetID(keys ...interface{}) {
	o.ID = keys[0].(string)
}

func (o *JournalLine) GetID(dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{o.ID}
}

func (o *JournalLine) PreSave(dbflex.IConnection) error {
	if o.ID == "" {
		o.ID = primitive.NewObjectID().Hex()
	}
	if o.Created.IsZero() {
		o.Created = time.Now()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *JournalLine) PostSave(dbflex.IConnection) error {
	return nil
}
func (o *JournalLine) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{}
}
