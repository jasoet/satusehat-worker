package entity

import (
	"encoding/json"
	shared "github.com/jasoet/fhir-worker/shared/model"
	"time"
)

type SatuSehatInternal struct {
	VisitID                   string           `db:"visit_id"`
	VisitDate                 time.Time        `db:"visit_date"`
	SatusehatPatientID        string           `db:"satusehat_patient_id"`
	VisitDetailJson           json.RawMessage  `db:"visit_detail"`
	VitalSignJson             json.RawMessage  `db:"vital_sign"`
	DiagnosisJsonArr          *json.RawMessage `db:"diagnosis"`           //Json Array
	LabJsonArr                *json.RawMessage `db:"lab"`                 //Json Array
	RadiologyJsonArr          *json.RawMessage `db:"radiology"`           //Json Array
	MedicationRequestJsonArr  *json.RawMessage `db:"medication_request"`  //Json Array
	MedicationDispenseJsonArr *json.RawMessage `db:"medication_dispense"` //Json Array
	ProcedureJsonArr          *json.RawMessage `db:"medical_procedure"`   //Json Array
	PublishDate               *time.Time       `db:"publish_date"`
	PublishRequest            *string          `db:"publish_request"`
	PublishResponse           *string          `db:"publish_response"`
	MappingErrors             *string          `db:"mapping_errors"`
	MappingStatus             MappingStatus    `db:"mapping_status"`
	PublishStatus             PublishStatus    `db:"publish_status"`
}

func (s *SatuSehatInternal) VisitDetail() *shared.VisitDetail {
	var visitDetail shared.VisitDetail
	err := json.Unmarshal(s.VisitDetailJson, &visitDetail)
	if err != nil {
		return nil
	}
	return &visitDetail
}

func (s *SatuSehatInternal) VitalSign() *shared.VitalSign {
	var vitalSign shared.VitalSign
	err := json.Unmarshal(s.VitalSignJson, &vitalSign)
	if err != nil {
		return nil
	}
	return &vitalSign
}

func (s *SatuSehatInternal) Diagnosis() *shared.DiagnosisList {
	if s.DiagnosisJsonArr == nil {
		return nil
	}

	var o shared.DiagnosisList
	err := json.Unmarshal(*s.DiagnosisJsonArr, &o)
	if err != nil || len(o) == 0 {
		return nil
	}
	return &o
}

func (s *SatuSehatInternal) Lab() *shared.ObservationLabList {
	if s.LabJsonArr == nil {
		return nil
	}
	var o shared.ObservationLabList
	err := json.Unmarshal(*s.LabJsonArr, &o)
	if err != nil || len(o) == 0 {
		return nil
	}
	return &o
}

func (s *SatuSehatInternal) Radiology() *shared.ObservationRadiologyList {
	if s.RadiologyJsonArr == nil {
		return nil
	}
	var o shared.ObservationRadiologyList
	err := json.Unmarshal(*s.RadiologyJsonArr, &o)
	if err != nil || len(o) == 0 {
		return nil
	}
	return &o
}

func (s *SatuSehatInternal) MedicationDispense() *shared.MedicationDispenseList {
	if s.MedicationRequestJsonArr == nil {
		return nil
	}
	var o shared.MedicationDispenseList
	err := json.Unmarshal(*s.MedicationDispenseJsonArr, &o)
	if err != nil || len(o) == 0 {
		return nil
	}
	return &o
}

func (s *SatuSehatInternal) MedicationRequest() *shared.MedicationRequestList {
	if s.MedicationRequestJsonArr == nil {
		return nil
	}
	var o shared.MedicationRequestList
	err := json.Unmarshal(*s.MedicationRequestJsonArr, &o)
	if err != nil || len(o) == 0 {
		return nil
	}
	return &o
}

func (s *SatuSehatInternal) Procedure() *shared.ProcedureList {
	if s.ProcedureJsonArr == nil {
		return nil
	}
	var o shared.ProcedureList
	err := json.Unmarshal(*s.ProcedureJsonArr, &o)
	if err != nil || len(o) == 0 {
		return nil
	}
	return &o
}
