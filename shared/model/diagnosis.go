package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Diagnosis struct {
	VisitID       string    `json:"visit_id" validate:"required"`
	DiagnosisCode string    `json:"diagnosis_code" validate:"required"`
	DiagnosisName string    `json:"diagnosis_name"  validate:"required"`
	DiagnosisDate time.Time `json:"diagnosis_date"  validate:"required"`
}

func (o *Diagnosis) Invalid() bool {
	val := validator.New()
	err := val.Struct(o)
	if err != nil {
		return true
	}
	return false
}
