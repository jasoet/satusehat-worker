package model

import "github.com/go-playground/validator/v10"

type Procedure struct {
	VisitId       int    `db:"visit_id" json:"visit_id" validate:"required"`
	ProcedureCode string `db:"procedure_code" json:"procedure_code" validate:"required"`
	ProcedureName string `db:"procedure_name" json:"procedure_name" validate:"required"`
}

func (o *Procedure) Invalid() bool {
	val := validator.New()
	err := val.Struct(o)
	if err != nil {
		return true
	}
	return false
}
