package resource

import (
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type ConditionDiagnosis struct {
	ConditionId        string `validate:"required"`
	EncounterId        string `validate:"required"`
	PatientSatuSehatId string `validate:"required"`
	PatientName        string `validate:"required"`
	Time               string `validate:"required"`
	IcdCode            string `validate:"required"`
	IcdName            string `validate:"required"`
}

func (o *ConditionDiagnosis) BundleEntry() (*fhir.BundleEntry, error) {
	return BundleEntry(o.Resource(), o.ConditionId, "ConditionDiagnosis")
}

func (o *ConditionDiagnosis) Resource() *fhir.Condition {
	condition := &fhir.Condition{
		ClinicalStatus: &fhir.CodeableConcept{
			Coding: []fhir.Coding{
				{
					System:  util.StrPtr("http://terminology.hl7.org/CodeSystem/condition-clinical"),
					Code:    util.StrPtr("active"),
					Display: util.StrPtr("Active"),
				},
			},
		},
		Subject: fhir.Reference{
			Reference: util.StrPtrFmt("Patient/%s", o.PatientSatuSehatId),
			Display:   util.StrPtr(o.PatientName),
		},
		Encounter: &fhir.Reference{
			Reference: util.StrPtrFmt("Encounter/%s", o.EncounterId),
			Display:   util.StrPtrFmt("Kunjungan %s. Di tanggal %s", o.PatientName, o.Time),
		},
		Category: []fhir.CodeableConcept{
			{
				Coding: []fhir.Coding{
					{
						System:  util.StrPtr("http://terminology.hl7.org/CodeSystem/condition-category"),
						Code:    util.StrPtr("encounter-diagnosis"),
						Display: util.StrPtr("Encounter Diagnosis"),
					},
				},
			},
		},
		Code: &fhir.CodeableConcept{
			Coding: []fhir.Coding{
				{
					System:  util.StrPtr("http://hl7.org/fhir/sid/icd-10"),
					Code:    util.StrPtr(o.IcdCode),
					Display: util.StrPtr(o.IcdName),
				},
			},
		},
	}

	return condition

}
