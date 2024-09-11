package model

type PatientType string

const (
	Outpatient PatientType = "Outpatient"
	Inpatient  PatientType = "Inpatient"
)

type MedicineType string

const (
	NonCompound MedicineType = "NonCompound"
	Compound    MedicineType = "Compound"
)
