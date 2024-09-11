//go:build !sleman

package simrs

import (
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/jasoet/fhir-worker/shared/model"
	"time"
)

func BuildVisit(m map[string]any) entity.Visit {
	v := entity.Visit{}

	v.VisitID = util.GetMapValueString(m, "visit_id", int64(0))
	v.PatientSatusehatID = util.GetMapValue(m, "patient_satusehat_id", "")
	v.PatientNIK = util.GetMapValueAsString(m, "patient_nik", "")
	v.PatientName = util.GetMapValueAsString(m, "patient_name", "")
	v.PatientSex = util.GetMapValueAsString(m, "patient_sex", "")

	v.PatientBirthDate = util.GetMapNullableValue[time.Time](m, "patient_birth_date")

	v.PatientAddress = util.GetMapValueAsString(m, "address", "")
	v.PractitionerNIK = util.GetMapValueAsString(m, "practitioner_nik", "")
	v.PractitionerSatusehatID = util.GetMapValueAsString(m, "practitioner_satusehat_id", "")
	v.PractitionerName = util.GetMapValueAsString(m, "practitioner_name", "")
	v.ClinicSatusehatID = util.GetMapValueAsString(m, "clinic_location_id", "")
	v.ClinicName = util.GetMapValueAsString(m, "clinic_name", "")

	v.Systole = util.GetMapValueString[int64](m, "visit_sistole", 0)
	v.Diastole = util.GetMapValueString[int64](m, "visit_diastole", 0)
	v.HeartRate = util.GetMapValueString[int64](m, "visit_heart_rate", 0)
	v.RespirationRate = util.GetMapValueString[int64](m, "visit_respiration_rate", 0)
	v.Temperature = util.GetMapValueAsString(m, "visit_temperature", "")
	v.OxygenSaturation = util.GetMapValueString[int64](m, "visit_spo2", 0)

	visitDate := util.GetMapValue(m, "visit_date", time.Time{})
	v.PeriodStartDate = &visitDate
	v.PeriodEndDate = &visitDate

	arrivedTime := util.GetMapNullableValue[time.Time](m, "registration_date")
	v.ArrivedStartTime = arrivedTime
	v.ArrivedEndTime = arrivedTime

	finishTime := util.GetMapNullableValue[time.Time](m, "examination_end_date")
	v.FinishStartTime = finishTime
	v.FinishEndTime = finishTime

	inProgressDate := util.GetMapNullableValue[time.Time](m, "examination_start_date")
	v.InProgressStartTime = inProgressDate
	v.InProgressEndTime = inProgressDate

	return v
}
