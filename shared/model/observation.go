package model

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type ObservationLab struct {
	VisitId          int              `db:"visit_id"`
	LabName          string           `db:"lab_name"`
	LabParameter     *json.RawMessage `db:"lab_parameter"`
	LabUnit          *json.RawMessage `db:"lab_unit"`
	LabNormal        *json.RawMessage `db:"lab_normal"`
	LabResult        *json.RawMessage `db:"lab_result"`
	LabFlag          *json.RawMessage `db:"lab_flag"`
	LabMethod        *json.RawMessage `db:"lab_method"`
	LabLoincCode     *json.RawMessage `db:"lab_loinc_code" validate:"required"`
	LabLoincName     *json.RawMessage `db:"lab_loinc_name" validate:"required"`
	PractitionerId   *string          `db:"practitioner_id"` //  validate:"required"`
	PractitionerName string           `db:"practitioner_name" validate:"required"`
}

func (o *ObservationLab) Invalid() bool {
	val := validator.New()
	err := val.Struct(o)
	if err != nil {
		return true
	}
	return false
}

type ObservationRadiology struct {
	VisitId          int              `db:"visit_id"`
	LabName          string           `db:"lab_name"`
	LabParameter     *json.RawMessage `db:"lab_parameter"`
	LabUnit          *json.RawMessage `db:"lab_unit"`
	LabNormal        *json.RawMessage `db:"lab_normal"`
	LabResult        *json.RawMessage `db:"lab_result"`
	LabFlag          *json.RawMessage `db:"lab_flag"`
	LabMethod        *json.RawMessage `db:"lab_method"`
	LabLoincCode     *json.RawMessage `db:"lab_loinc_code" validate:"required"`
	LabLoincName     *json.RawMessage `db:"lab_loinc_name" validate:"required"`
	PractitionerId   *string          `db:"practitioner_id"` //  validate:"required"`
	PractitionerName string           `db:"practitioner_name" validate:"required"`
}

func (o *ObservationRadiology) Invalid() bool {
	val := validator.New()
	err := val.Struct(o)
	if err != nil {
		return true
	}
	return false
}
