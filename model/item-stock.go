package model

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemStock struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id" key:"1" form_read_only_edit:"1" form_section:"General" form_section_auto_col:"2"`
	ItemID            string `form_required:"1" form_section:"General"`
	StockType         string `form_required:"1" form_section:"General"`
	Stock             float64
	Created           time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info" form_section_auto_col:"2"`
	LastUpdate        time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info"`
}

func (o *ItemStock) TableName() string {
	return "ItemStocks"
}

func (o *ItemStock) FK() []*orm.FKConfig {
	return orm.DefaultRelationManager().FKs(o)
}

func (o *ItemStock) ReverseFK() []*orm.ReverseFKConfig {
	return orm.DefaultRelationManager().ReverseFKs(o)
}

func (o *ItemStock) SetID(keys ...interface{}) {
	o.ID = keys[0].(string)
}

func (o *ItemStock) GetID(dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{o.ID}
}

func (o *ItemStock) PreSave(dbflex.IConnection) error {
	if o.ID == "" {
		o.ID = primitive.NewObjectID().Hex()
	}
	if o.Created.IsZero() {
		o.Created = time.Now()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *ItemStock) PostSave(dbflex.IConnection) error {
	return nil
}
func (o *ItemStock) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{
		{Name: "idx_item_id", Fields: []string{"ItemID"}},
		{Name: "idx_stock_type", Fields: []string{"StockType"}},
	}
}
