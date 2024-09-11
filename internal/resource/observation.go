package resource

import (
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type Observation struct {
	ObservationId           string
	EncounterId             string
	PatientSatuSehatId      string
	PatientName             string
	Time                    string
	LoincCode               string
	LoincDisplay            string
	PractitionerSatuSehatId string
	PractitionerName        string
	ValueQuantity           *ObservationValue
	ValueCode               *ObservationValueCode
}

type ObservationValue struct {
	Value string
	Unit  string
	Code  string
}

type ObservationValueCode struct {
	Code    string
	Display string
}

func (o *Observation) BundleEntry() (*fhir.BundleEntry, error) {
	return BundleEntry(o.Resource(), o.ObservationId, "Observation")
}

func (o *Observation) Resource() fhir.Observation {
	observation := fhir.Observation{
		Status: fhir.ObservationStatusFinal,
		Subject: &fhir.Reference{
			Reference: util.StrPtrFmt("Patient/%s", o.PatientSatuSehatId),
			Display:   util.StrPtr(o.PatientName),
		},
		Encounter: &fhir.Reference{
			Reference: util.StrPtrFmt("Encounter/%s", o.EncounterId),
			Display:   util.StrPtrFmt("Kunjungan %s. Di tanggal %s", o.PatientName, o.Time),
		},
		Performer: []fhir.Reference{
			{
				Reference: util.StrPtrFmt("Practitioner/%s", o.PractitionerSatuSehatId),
				Display:   util.StrPtr(o.PractitionerName),
			},
		},
		EffectiveDateTime: util.StrPtr(o.Time),
		Code: fhir.CodeableConcept{
			Coding: []fhir.Coding{
				{
					System:  util.StrPtr("http://loinc.org"),
					Code:    util.StrPtr(o.LoincCode),
					Display: util.StrPtr(o.LoincDisplay),
				},
			},
		},
	}

	if o.ValueQuantity != nil {
		observation.ValueQuantity = &fhir.Quantity{
			System: util.StrPtr("http://unitsofmeasure.org"),
			Value:  util.JsonNumber(o.ValueQuantity.Value),
			Unit:   util.StrPtr(o.ValueQuantity.Unit),
			Code:   util.StrPtr(o.ValueQuantity.Code),
		}
	}

	if o.ValueCode != nil {
		observation.ValueCodeableConcept = &fhir.CodeableConcept{
			Coding: []fhir.Coding{
				{
					System:  util.StrPtr("http://loinc.org"),
					Code:    util.StrPtr(o.ValueCode.Code),
					Display: util.StrPtr(o.ValueCode.Display),
				},
			},
		}
	}

	return observation

}
