package resource

import (
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/jasoet/fhir-worker/shared/model"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type MedicationRequest struct {
	MedicationId        string             `validate:"required"`
	MedicationRequestId string             `validate:"required"`
	EncounterId         string             `validate:"required"`
	OrganizationId      string             `validate:"required"`
	PrescriptionId      string             `validate:"required"`
	KfaCode             string             `validate:"required"`
	KfaDisplay          string             `validate:"required"`
	Type                model.MedicineType `validate:"required"`
	PatientType         model.PatientType  `validate:"required"`
	PatientId           string             `validate:"required"`
	PatientName         string             `validate:"required"`
	PractitionerId      string             `validate:"required"`
	PractitionerName    string             `validate:"required"`
	Date                string             `validate:"required"`
}

func (o *MedicationRequest) PatientTypeCoding() fhir.Coding {
	var coding fhir.Coding
	if o.PatientType == model.Outpatient {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.hl7.org/CodeSystem/medicationrequest-category"),
			Code:    util.StrPtr("outpatient"),
			Display: util.StrPtr("Outpatient"),
		}
	} else {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.hl7.org/CodeSystem/medicationrequest-category"),
			Code:    util.StrPtr("inpatient"),
			Display: util.StrPtr("Inpatient"),
		}
	}

	return coding

}

func (o *MedicationRequest) UsageCoding() fhir.Coding {
	var coding fhir.Coding
	if o.Type == model.NonCompound {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.kemkes.go.id/CodeSystem/medication-type"),
			Code:    util.StrPtr("NC"),
			Display: util.StrPtr("Non-compound"),
		}
	} else {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.kemkes.go.id/CodeSystem/medication-type"),
			Code:    util.StrPtr("C"),
			Display: util.StrPtr("Compound"),
		}
	}
	return coding
}

func (o *MedicationRequest) Resources() (*fhir.Medication, *fhir.MedicationRequest) {
	official := fhir.IdentifierUseOfficial
	priority := fhir.RequestPriorityRoutine

	medication := &fhir.Medication{
		Identifier: []fhir.Identifier{
			{
				System: util.StrPtrFmt("http://sys-ids.kemkes.go.id/medication/%s", o.OrganizationId),
				Use:    &official,
				Value:  util.StrPtr(o.PrescriptionId),
			},
		},
		Code: &fhir.CodeableConcept{
			Coding: []fhir.Coding{
				{
					System:  util.StrPtr("http://sys-ids.kemkes.go.id/kfa"),
					Code:    util.StrPtr(o.KfaCode),
					Display: util.StrPtr(o.KfaDisplay),
				},
			},
		},
		Status: util.StrPtr("active"),
		Extension: []fhir.Extension{
			{
				Url: "https://fhir.kemkes.go.id/r4/StructureDefinition/MedicationType",
				ValueCodeableConcept: &fhir.CodeableConcept{
					Coding: []fhir.Coding{
						o.UsageCoding(),
					},
				},
			},
		},
	}

	medicationRequest := &fhir.MedicationRequest{
		Identifier: []fhir.Identifier{
			{
				System: util.StrPtrFmt("http://sys-ids.kemkes.go.id/prescription/%s", o.OrganizationId),
				Use:    &official,
				Value:  util.StrPtr(o.PrescriptionId),
			},
		},
		Status: "completed",
		Intent: "order",
		Category: []fhir.CodeableConcept{
			{
				Coding: []fhir.Coding{
					o.PatientTypeCoding(),
				},
			},
		},
		Priority: &priority,
		MedicationReference: fhir.Reference{
			Reference: util.StrPtrFmt("urn:uuid:%s", o.MedicationId),
			Display:   util.StrPtr(o.KfaDisplay),
		},
		Subject: fhir.Reference{
			Reference: util.StrPtrFmt("Patient/%s", o.PatientId),
			Display:   util.StrPtr(o.PatientName),
		},
		DispenseRequest: &fhir.MedicationRequestDispenseRequest{
			Performer: &fhir.Reference{
				Reference: util.StrPtrFmt("Organization/%s", o.OrganizationId),
			},
		},
		Encounter: &fhir.Reference{
			Reference: util.StrPtrFmt("Encounter/%s", o.EncounterId),
		},
		AuthoredOn: util.StrPtr(o.Date),
		Requester: &fhir.Reference{
			Reference: util.StrPtrFmt("Practitioner/%s", o.PractitionerId),
			Display:   util.StrPtr(o.PatientName),
		},
	}

	return medication, medicationRequest
}

func (o *MedicationRequest) BundleEntries() ([]fhir.BundleEntry, error) {
	var result []fhir.BundleEntry
	medication, medicationRequest := o.Resources()

	medicationEntry, err := BundleEntry(medication, o.MedicationId, "Medication")
	if err != nil {
		return nil, err
	}

	result = append(result, *medicationEntry)

	medicationRequestEntry, err := BundleEntry(medicationRequest, o.MedicationRequestId, "MedicationRequest", WithRemoveKey("medicationCodeableConcept"))
	if err != nil {
		return nil, err
	}
	result = append(result, *medicationRequestEntry)
	return result, nil
}
