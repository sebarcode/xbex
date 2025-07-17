package logic

import (
	"errors"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/datahub"
	"github.com/sebarcode/xbex/model"
)

func addJournalLog(db *datahub.Hub, journalID, logType, logMessage string) {
	log := &model.JournalLog{
		JournalID:  journalID,
		LogType:    logType,
		LogMessage: logMessage,
	}
	_ = db.Insert(log)
}

type JournalHandler struct {
}

type ApproveRequest struct {
	Op            string
	JournalID     string
	Justification string
}

func (obj *JournalHandler) Approve(ctx *kaos.Context, payload *ApproveRequest) (*model.JournalHeader, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return nil, errors.New("missingDBConn")
	}

	header, err := datahub.GetByID(h, new(model.JournalHeader), payload.JournalID)
	if err != nil {
		return nil, errors.New("journalNotFound")
	}
	if header.Status != "Submitted" {
		return nil, errors.New("journalNotSubmitted")
	}
	if payload.Op == "Reject" && payload.Justification == "" {
		return nil, errors.New("justificationRequired")
	}
	header.Status = payload.Op
	header.Justification = payload.Justification
	err = h.Update(header)
	if err != nil {
		return nil, ctx.Log().Error2(err.Error(), "failedToUpdateJournal")
	}
	addJournalLog(h, payload.JournalID, header.Status, payload.Justification)

	return header, nil
}

func (obj *JournalHandler) Post(ctx *kaos.Context, journalID string) (*model.JournalHeader, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return nil, errors.New("missingDBConn")
	}

	header, err := datahub.GetByID(h, new(model.JournalHeader), journalID)
	if err != nil {
		return nil, errors.New("journalNotFound")
	}
	if header.Status != "Approved" {
		return nil, errors.New("journalNotApproved")
	}
	header.Status = "Posted"
	err = h.Update(header)
	if err != nil {
		return nil, ctx.Log().Error2(err.Error(), "failedToUpdateJournal")
	}
	addJournalLog(h, journalID, header.Status, "Journal posted")
	return header, nil
}

func (obj *JournalHandler) Cancel(ctx *kaos.Context, journalID, justification string) (*model.JournalHeader, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return nil, errors.New("missingDBConn")
	}

	header, err := datahub.GetByID(h, new(model.JournalHeader), journalID)
	if err != nil {
		return nil, errors.New("journalNotFound")
	}
	if header.Status == "Posted" {
		return nil, errors.New("journalAlreadyPosted")
	}
	if header.Status == "Cancelled" {
		return nil, errors.New("journalAlreadyCancelled")
	}
	if justification == "" {
		return nil, errors.New("justificationRequired")
	}
	header.Status = "Cancelled"
	header.Justification = justification
	err = h.Update(header)
	if err != nil {
		return nil, ctx.Log().Error2(err.Error(), "failedToUpdateJournal")
	}
	addJournalLog(h, journalID, header.Status, justification)
	return header, nil
}

func (obj *JournalHandler) Copy(ctx *kaos.Context, journalID string) (*model.JournalHeader, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return nil, errors.New("missingDBConn")
	}

	header, err := datahub.GetByID(h, new(model.JournalHeader), journalID)
	if err != nil {
		return nil, errors.New("journalNotFound")
	}

	return func() (*model.JournalHeader, error) {
		dbTx, err := h.BeginTx()
		if err != nil {
			return nil, ctx.Log().Error2(err.Error(), "failedToBeginTransaction")
		}
		newHeader := *header
		newHeader.ID = ""
		newHeader.Status = "Draft"
		newHeader.Justification = ""
		err = dbTx.Insert(&newHeader)
		if err != nil {
			dbTx.Rollback()
			return nil, ctx.Log().Error2(err.Error(), "failedToCopyJournal")
		}

		lines, err := datahub.FindByFilter(h, new(model.JournalLine), dbflex.Eq("JournalID", journalID))
		if err != nil {
			dbTx.Rollback()
			return nil, ctx.Log().Error2(err.Error(), "failedToGetJournalLines")
		}
		for _, line := range lines {
			newLine := *line
			newLine.ID = ""
			newLine.JournalID = newHeader.ID
			err = dbTx.Insert(&newLine)
			if err != nil {
				dbTx.Rollback()
				return nil, ctx.Log().Error2(err.Error(), "failedToCopyJournalLine")
			}
		}

		addJournalLog(dbTx, newHeader.ID, "Copy", "Journal copied from "+journalID)

		return &newHeader, dbTx.Commit()
	}()
}

func (obj *JournalHandler) Submit(ctx *kaos.Context, journalID string) (*model.JournalHeader, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return nil, errors.New("missingDBConn")
	}

	header, err := datahub.GetByID(h, new(model.JournalHeader), journalID)
	if err != nil {
		return nil, errors.New("journalNotFound")
	}

	if header.Status != "Draft" {
		return nil, errors.New("journalNotDraft")
	}

	lines, err := datahub.FindByFilter(h, new(model.JournalLine), dbflex.Eq("JournalID", journalID))
	if err != nil {
		return nil, ctx.Log().Error2(err.Error(), "failedToGetJournalLines")
	}
	if len(lines) == 0 {
		return nil, errors.New("journalNoLines")
	}

	header.Status = "Submitted"
	err = h.Update(header)
	if err != nil {
		return nil, ctx.Log().Error2(err.Error(), "failedToUpdateJournal")
	}
	addJournalLog(h, journalID, header.Status, "Journal submitted")
	return header, nil
}
