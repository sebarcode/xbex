package model

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JournalLog struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id" key:"1" form_read_only_edit:"1" form_section:"General" form_section_auto_col:"2"`
	JournalID         string
	LogType           string
	LogMessage        string
	UserID            string
	Created           time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info" form_section_auto_col:"2"`
	LastUpdate        time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info"`
}

func (o *JournalLog) TableName() string {
	return "JournalLogs"
}

func (o *JournalLog) FK() []*orm.FKConfig {
	return orm.DefaultRelationManager().FKs(o)
}

func (o *JournalLog) ReverseFK() []*orm.ReverseFKConfig {
	return orm.DefaultRelationManager().ReverseFKs(o)
}

func (o *JournalLog) SetID(keys ...interface{}) {
	o.ID = keys[0].(string)
}

func (o *JournalLog) GetID(dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{o.ID}
}

func (o *JournalLog) PreSave(dbflex.IConnection) error {
	if o.ID == "" {
		o.ID = primitive.NewObjectID().Hex()
	}
	if o.Created.IsZero() {
		o.Created = time.Now()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *JournalLog) PostSave(dbflex.IConnection) error {
	return nil
}
func (o *JournalLog) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{
		{Name: "idx_journal_log_journal_id", Fields: []string{"JournalID", "LogType"}},
		{Name: "idx_journal_log_user_id", Fields: []string{"UserID"}},
		{Name: "idx_journal_log_logType", Fields: []string{"LogType"}},
	}
}
