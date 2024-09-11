package resource

import (
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type EncounterDiagnosis struct {
	Id      string
	Display string
}

type Encounter struct {
	EncounterId             string `validate:"required"`
	PatientSatuSehatId      string `validate:"required"`
	PatientName             string `validate:"required"`
	PractitionerSatuSehatId string `validate:"required"`
	PractitionerName        string `validate:"required"`
	OrganizationId          string `validate:"required"`
	LocationName            string `validate:"required"`
	LocationId              string `validate:"required"`
	PeriodStartDate         string `validate:"required"`
	PeriodEndDate           string `validate:"required"`
	ArrivedStartTime        string `validate:"required"`
	ArrivedEndTime          string `validate:"required"`
	InProgressStartTime     string `validate:"required"`
	InProgressEndTime       string `validate:"required"`
	FinishStartTime         string `validate:"required"`
	FinishEndTime           string `validate:"required"`
	Diagnosis               []EncounterDiagnosis
}

func (o *Encounter) BundleEntry() (*fhir.BundleEntry, error) {
	return BundleEntry(o.Resource(), o.EncounterId, "Encounter")
}

func (o *Encounter) Resource() fhir.Encounter {
	encounter := fhir.Encounter{
		Identifier: []fhir.Identifier{
			{
				System: util.StrPtrFmt("http://sys-ids.kemkes.go.id/encounter/%s", o.OrganizationId),
				Value:  util.StrPtr(o.PatientSatuSehatId),
			},
		},
		Status: fhir.EncounterStatusFinished,
		Subject: &fhir.Reference{
			Reference: util.StrPtrFmt("Patient/%s", o.PatientSatuSehatId),
			Display:   util.StrPtr(o.PatientName),
		},
		Participant: []fhir.EncounterParticipant{
			{
				Type: []fhir.CodeableConcept{
					{
						Coding: []fhir.Coding{
							{
								System:  util.StrPtr("http://terminology.hl7.org/CodeSystem/v3-ParticipationType"),
								Code:    util.StrPtr("ATND"),
								Display: util.StrPtr("attender"),
							},
						},
					},
				},
				Individual: &fhir.Reference{
					Reference: util.StrPtrFmt("Practitioner/%s", o.PractitionerSatuSehatId),
					Display:   util.StrPtr(o.PractitionerName),
				},
			},
		},
		Period: &fhir.Period{
			Start: util.StrPtr(o.PeriodStartDate),
			End:   util.StrPtr(o.PeriodEndDate),
		},
		Class: fhir.Coding{
			System:  util.StrPtr("http://terminology.hl7.org/CodeSystem/v3-ActCode"),
			Code:    util.StrPtr("AMB"),
			Display: util.StrPtr("ambulatory"),
		},
		ServiceProvider: &fhir.Reference{
			Reference: util.StrPtrFmt("Organization/%s", o.OrganizationId),
		},
		Location: []fhir.EncounterLocation{
			{
				Location: fhir.Reference{
					Reference: util.StrPtrFmt("Location/%s", o.LocationId),
					Display:   util.StrPtr(o.LocationName),
				},
			},
		},
		StatusHistory: []fhir.EncounterStatusHistory{
			{
				Status: fhir.EncounterStatusArrived,
				Period: fhir.Period{
					Start: util.StrPtr(o.ArrivedStartTime),
					End:   util.StrPtr(o.ArrivedEndTime),
				},
			},
			{
				Status: fhir.EncounterStatusInProgress,
				Period: fhir.Period{
					Start: util.StrPtr(o.InProgressStartTime),
					End:   util.StrPtr(o.InProgressEndTime),
				},
			},
			{
				Status: fhir.EncounterStatusFinished,
				Period: fhir.Period{
					Start: util.StrPtr(o.FinishStartTime),
					End:   util.StrPtr(o.FinishEndTime),
				},
			},
		},
	}

	var encounterDiagnosis []fhir.EncounterDiagnosis

	for _, diagnosis := range o.Diagnosis {
		entry := fhir.EncounterDiagnosis{
			Condition: fhir.Reference{
				Reference: util.StrPtrFmt("urn:uuid:%s", diagnosis.Id),
				Display:   util.StrPtr(diagnosis.Display),
			},
			Use: &fhir.CodeableConcept{
				Coding: []fhir.Coding{
					{
						System:  util.StrPtr("http://terminology.hl7.org/CodeSystem/diagnosis-role"),
						Code:    util.StrPtr("DD"),
						Display: util.StrPtr("Discharge diagnosis"),
					},
				},
			},
		}

		encounterDiagnosis = append(encounterDiagnosis, entry)
	}

	encounter.Diagnosis = encounterDiagnosis

	return encounter

}
