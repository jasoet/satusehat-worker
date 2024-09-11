package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type MedicationRequest struct {
	VisitId          int          `json:"visit_id"`
	PatientType      PatientType  `json:"patient_type"  validate:"required"`
	Date             *time.Time   `json:"date"  validate:"required"`
	MedicineCode     *string      `json:"medicine_code"`
	PrescriptionId   int          `json:"prescription_id" validate:"required"`
	KfaCode          *string      `json:"kfa_code"`
	KfaName          *string      `json:"kfa_name"`
	Type             MedicineType `json:"type"  validate:"required"`
	PractitionerId   *string      `json:"practitioner_id"  validate:"required"`
	PractitionerName *string      `json:"practitioner_name"  validate:"required"`
	Amount           float64      `json:"amount"`
	Unit             string       `json:"unit"`
}

func (o *MedicationRequest) Invalid() bool {
	val := validator.New()
	err := val.Struct(o)
	if err != nil {
		return true
	}
	return false
}

type MedicationDispense struct {
	VisitId               int          `json:"visit_id"`
	PatientType           PatientType  `json:"patient_type" validate:"required"`
	Date                  *time.Time   `json:"date" validate:"required"`
	MedicineCode          string       `json:"medicine_code"`
	PrescriptionId        int          `json:"prescription_id" validate:"required"`
	KfaCode               *string      `json:"kfa_code"`
	KfaName               *string      `json:"kfa_name"`
	Type                  MedicineType `json:"type" validate:"required"`
	PractitionerId        *string      `json:"practitioner_id" validate:"required"`
	PractitionerName      *string      `json:"practitioner_name" validate:"required"`
	BatchNumber           string       `json:"batch_number" validate:"required"`
	ExpiredDate           *time.Time   `json:"expired_date" validate:"required"`
	PrescriptionStartDate *time.Time   `json:"prescription_start_date" validate:"required"`
	HandoverDate          *time.Time   `json:"drug_received_date" validate:"required"`
}

func (o *MedicationDispense) Invalid() bool {
	val := validator.New()
	err := val.Struct(o)
	if err != nil {
		return true
	}
	return false
}
