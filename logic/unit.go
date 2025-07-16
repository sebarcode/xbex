package logic

import (
	"git.kanosolution.net/kano/dbflex"
	"github.com/ariefdarmawan/datahub"
	"github.com/sebarcode/xbex/model"
)

func CalcUnit(db *datahub.Hub, fromQty, toQty float64, fromUnitID string, toUnitID string) float64 {
	if fromQty == 0 {
		return 0
	}

	if fromUnitID == toUnitID {
		return fromQty
	}

	reverse := false
	factor, err := datahub.GetByFilter(db, new(model.UnitFactor), dbflex.Eqs("FromUnitID", fromUnitID, "ToUnitID", toUnitID))
	if err != nil {
		factor, err = datahub.GetByFilter(db, new(model.UnitFactor), dbflex.Eqs("FromUnitID", toUnitID, "ToUnitID", fromUnitID))
		if err != nil {
			return fromQty
		}
		reverse = true
	}

	if reverse {
		if factor.FromQty == 0 {
			return fromQty
		}
		return fromQty * (factor.ToQty / factor.FromQty)
	}

	if factor.ToQty == 0 {
		return fromQty
	}
	return fromQty * (factor.FromQty / factor.ToQty)
}
