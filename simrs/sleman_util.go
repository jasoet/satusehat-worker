//go:build sleman

package simrs

import (
	"fmt"
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/jasoet/fhir-worker/shared/model"
	"strings"
	"time"
)

func BuildVisit(m map[string]any) model.Visit {
	v := model.Visit{}

	v.VisitID = util.GetMapValue(m, "visit_id", "")
	v.PatientSatusehatID = util.GetMapValue(m, "patient_satusehat_id", "")
	v.PatientNIK = util.GetMapValue(m, "patient_nik", "")
	v.PatientName = util.GetMapValue(m, "patient_name", "")
	gender := util.GetMapValue(m, "patient_sex", "")

	if gender == "1" {
		v.PatientSex = "male"
	} else {
		v.PatientSex = "female"
	}

	patientBirthDate := util.GetMapValue[time.Time](m, "patient_birth_date", time.Time{})
	v.PatientBirthDate = &patientBirthDate

	v.PatientAddress = util.GetMapValue(m, "patient_address", "")
	v.PractitionerNIK = util.GetMapValue(m, "practitioner_nik", "")
	v.PractitionerSatusehatID = util.GetMapValue(m, "practitioner_satusehat_id", "")
	v.PractitionerName = util.GetMapValue(m, "practitioner_name", "")
	v.ClinicSatusehatID = util.GetMapValue(m, "clinic_satusehat_id", "")
	v.ClinicName = util.GetMapValue(m, "clinic_name", "")

	bloodPressure := util.GetMapValue(m, "blood_pressure", "")

	v.Systole, v.Diastole = bloodPressureBreakdown(bloodPressure)

	v.HeartRate = util.GetMapValue(m, "heart_rate", "")
	v.RespirationRate = util.GetMapValue(m, "respiration_rate", "")
	v.Temperature = util.GetMapValue(m, "temperature", "")

	visitDate := util.GetMapValue(m, "visit_date", time.Time{})
	v.PeriodStartDate = visitDate
	v.PeriodEndDate = visitDate

	arrivedTime := util.GetMapNullableValue[time.Time](m, "visit_arrived_time")
	v.ArrivedStartTime = arrivedTime
	v.ArrivedEndTime = arrivedTime

	finishTime := util.GetMapNullableValue[time.Time](m, "visit_end_time")
	v.FinishStartTime = finishTime
	v.FinishEndTime = finishTime

	inProgressDate := util.GetMapNullableValue[time.Time](m, "visit_inprogress_date")
	inProgressHour := util.GetMapValue(m, "visit_inprogress_hour", "")

	if inProgressDate != nil && util.StringNotEmpty(inProgressHour) {
		dateTime, err := combineDateTime(*inProgressDate, inProgressHour)
		if err != nil {
			v.InProgressStartTime = inProgressDate
		} else {
			v.InProgressStartTime = dateTime
			v.InProgressEndTime = dateTime
		}
	}

	return v
}

func bloodPressureBreakdown(bloodPressure string) (string, string) {
	values := strings.Split(bloodPressure, "/")

	// Check if we have exactly two values
	if len(values) != 2 {
		return "", ""
	}

	// Trim any whitespace from the values
	systolic := strings.TrimSpace(values[0])
	diastolic := strings.TrimSpace(values[1])

	return systolic, diastolic
}

func combineDateTime(inProgressDate time.Time, inProgressHour string) (*time.Time, error) {
	// Parse the hour string
	hourMinute, err := time.Parse("15:04", inProgressHour)
	if err != nil {
		return nil, fmt.Errorf("invalid hour format: %v", err)
	}

	// Combine date and time
	combinedDateTime := time.Date(
		inProgressDate.Year(),
		inProgressDate.Month(),
		inProgressDate.Day(),
		hourMinute.Hour(),
		hourMinute.Minute(),
		0, // seconds
		0, // nanoseconds
		inProgressDate.Location(),
	)

	return &combinedDateTime, nil
}

func BuildDiagnosis(m map[string]any) model.Diagnosis {
	o := model.Diagnosis{}

	o.VisitID = util.GetMapValue(m, "visit_id", "")
	o.DiagnosisDate = util.GetMapValue[time.Time](m, "diagnosis_date", time.Time{})
	o.DiagnosisCode = util.GetMapValue(m, "diagnosis_code", "")
	o.DiagnosisName = util.GetMapValue(m, "diagnosis_name", "")

	return o
}

func BuildMedicationRequest(m map[string]any) model.MedicationRequest {
	o := model.MedicationRequest{}

	return o
}
