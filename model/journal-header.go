package model

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JournalHeader struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id" key:"1" form_read_only_edit:"1" form_section:"General" form_section_auto_col:"2"`
	JournalType       string
	PostingProfile    string
	Name              string `form_required:"1" form_section:"General"`
	Status            string `form_required:"1" form_items:"Draft|Submitted|Approved|Rejected|Posted|Cancelled"`
	JournalDate       time.Time
	LineCount         int
	TotalQty          float64
	TotalAmount       float64
	Created           time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info" form_section_auto_col:"2"`
	LastUpdate        time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info"`
}

func (o *JournalHeader) TableName() string {
	return "JournalHeaders"
}

func (o *JournalHeader) FK() []*orm.FKConfig {
	return orm.DefaultRelationManager().FKs(o)
}

func (o *JournalHeader) ReverseFK() []*orm.ReverseFKConfig {
	return orm.DefaultRelationManager().ReverseFKs(o)
}

func (o *JournalHeader) SetID(keys ...interface{}) {
	o.ID = keys[0].(string)
}

func (o *JournalHeader) GetID(dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{o.ID}
}

func (o *JournalHeader) PreSave(dbflex.IConnection) error {
	if o.ID == "" {
		o.ID = primitive.NewObjectID().Hex()
	}
	if o.Created.IsZero() {
		o.Created = time.Now()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *JournalHeader) PostSave(dbflex.IConnection) error {
	return nil
}
func (o *JournalHeader) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{}
}
