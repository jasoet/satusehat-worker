package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Visit struct {
	VisitID                 string
	PatientSatusehatID      string
	PatientNIK              string
	PatientName             string
	PatientSex              string
	PatientBirthDate        *time.Time
	PatientAddress          string
	PractitionerNIK         string
	PractitionerSatusehatID string
	PractitionerName        string
	ClinicSatusehatID       string
	ClinicName              string
	Systole                 string
	Diastole                string
	HeartRate               string
	RespirationRate         string
	OxygenSaturation        string
	Temperature             string
	PeriodStartDate         time.Time
	PeriodEndDate           time.Time
	ArrivedStartTime        *time.Time
	ArrivedEndTime          *time.Time
	InProgressStartTime     *time.Time
	InProgressEndTime       *time.Time
	FinishStartTime         *time.Time
	FinishEndTime           *time.Time
}

type VisitDetail struct {
	VisitId             string     `json:"visit_id" validate:"required"`
	PatientSatusehatId  string     `json:"patient_satusehat_id" validate:"required"`
	PatientNik          string     `json:"patient_nik" `
	PatientName         string     `json:"patient_name" validate:"required"`
	PatientSex          string     `json:"patient_sex"`
	PatientBirthDate    *time.Time `json:"patient_birth_date"`
	PatientAddress      string     `json:"patient_address"`
	PractitionerNik     string     `json:"practitioner_nik"`
	PractitionerId      string     `json:"practitioner_satusehat_id" validate:"required"`
	PractitionerName    string     `json:"practitioner_name" validate:"required"`
	ClinicName          string     `json:"clinic_name" validate:"required"`
	ClinicSatuSehatId   string     `json:"clinic_id" validate:"required"`
	PeriodStartDate     time.Time  `json:"period_start_date" validate:"required"`
	PeriodEndDate       time.Time  `json:"period_end_date" validate:"required"`
	ArrivedStartTime    *time.Time `json:"arrived_start_time"  validate:"required"`
	ArrivedEndTime      *time.Time `json:"arrived_end_time"  validate:"required"`
	InProgressStartTime *time.Time `json:"in_progress_start_time"  validate:"required"`
	InProgressEndTime   *time.Time `json:"in_progress_end_time"  validate:"required"`
	FinishStartTime     *time.Time `json:"finish_start_time"  validate:"required"`
	FinishEndTime       *time.Time `json:"finish_end_time"  validate:"required"`
}

func (v VisitDetail) Invalid() error {
	val := validator.New()
	return val.Struct(v)
}

type VitalSign struct {
	Systole          string `json:"sistole"`
	Diastole         string `json:"diastole"`
	HeartRate        string `json:"heart_rate"`
	RespirationRate  string `json:"respiration_rate"`
	Temperature      string `json:"temperature"`
	OxygenSaturation string `json:"oxygen_saturation"`
}

func (v *Visit) VitalSign() VitalSign {
	return VitalSign{
		Systole:          v.Systole,
		Diastole:         v.Diastole,
		HeartRate:        v.HeartRate,
		RespirationRate:  v.RespirationRate,
		Temperature:      v.Temperature,
		OxygenSaturation: v.OxygenSaturation,
	}
}

func (v *Visit) VisitDetail() VisitDetail {
	return VisitDetail{
		VisitId:             v.VisitID,
		PatientSatusehatId:  v.PatientSatusehatID,
		PatientNik:          v.PatientNIK,
		PatientName:         v.PatientName,
		PatientSex:          v.PatientSex,
		PatientBirthDate:    v.PatientBirthDate,
		PatientAddress:      v.PatientAddress,
		ClinicName:          v.ClinicName,
		ClinicSatuSehatId:   v.ClinicSatusehatID,
		PeriodStartDate:     v.PeriodStartDate,
		PeriodEndDate:       v.PeriodEndDate,
		PractitionerNik:     v.PractitionerNIK,
		PractitionerId:      v.PractitionerSatusehatID,
		PractitionerName:    v.PractitionerName,
		ArrivedStartTime:    v.ArrivedStartTime,
		ArrivedEndTime:      v.ArrivedEndTime,
		InProgressStartTime: v.InProgressStartTime,
		InProgressEndTime:   v.InProgressEndTime,
		FinishStartTime:     v.FinishStartTime,
		FinishEndTime:       v.FinishEndTime,
	}
}
