package resource

import (
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/jasoet/fhir-worker/shared/model"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type MedicationDispense struct {
	MedicationId         string             `validate:"required"`
	MedicationDispenseId string             `validate:"required"`
	EncounterId          string             `validate:"required"`
	OrganizationId       string             `validate:"required"`
	PrescriptionId       string             `validate:"required"`
	KfaCode              string             `validate:"required"`
	KfaDisplay           string             `validate:"required"`
	Type                 model.MedicineType `validate:"required"`
	PatientType          model.PatientType  `validate:"required"`
	PatientId            string             `validate:"required"`
	PatientName          string             `validate:"required"`
	PractitionerId       string             `validate:"required"`
	PractitionerName     string             `validate:"required"`
	PreparedDate         string             `validate:"required"`
	HandoverDate         string             `validate:"required"`
	BatchNumber          string             `validate:"required"`
	ExpirationDate       string             `validate:"required"`
}

func (o *MedicationDispense) PatientTypeCoding() fhir.Coding {
	var coding fhir.Coding
	if o.PatientType == model.Outpatient {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.hl7.org/fhir/CodeSystem/medicationdispense-category"),
			Code:    util.StrPtr("outpatient"),
			Display: util.StrPtr("Outpatient"),
		}
	} else if o.PatientType == model.Inpatient {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.hl7.org/fhir/CodeSystem/medicationdispense-category"),
			Code:    util.StrPtr("inpatient"),
			Display: util.StrPtr("Inpatient"),
		}
	}

	return coding

}
func (o *MedicationDispense) UsageCoding() fhir.Coding {
	var coding fhir.Coding
	if o.Type == model.NonCompound {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.kemkes.go.id/CodeSystem/medication-type"),
			Code:    util.StrPtr("NC"),
			Display: util.StrPtr("Non-compound"),
		}
	} else if o.Type == model.Compound {
		coding = fhir.Coding{
			System:  util.StrPtr("http://terminology.kemkes.go.id/CodeSystem/medication-type"),
			Code:    util.StrPtr("C"),
			Display: util.StrPtr("Compound"),
		}
	}
	return coding
}

func (o *MedicationDispense) Resources() (*fhir.Medication, *fhir.MedicationDispense) {
	official := fhir.IdentifierUseOfficial

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
		Batch: &fhir.MedicationBatch{
			LotNumber:      util.StrPtr(o.BatchNumber),
			ExpirationDate: util.StrPtr(o.ExpirationDate),
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

	medicationDispense := &fhir.MedicationDispense{
		Identifier: []fhir.Identifier{
			{
				System: util.StrPtrFmt("http://sys-ids.kemkes.go.id/prescription/%s", o.OrganizationId),
				Use:    &official,
				Value:  util.StrPtr(o.PrescriptionId),
			},
		},
		Performer: []fhir.MedicationDispensePerformer{
			{
				Actor: fhir.Reference{
					Reference: util.StrPtrFmt("Practitioner/%s", o.PractitionerId),
					Display:   util.StrPtr(o.PatientName),
				},
			},
		},
		Status: "completed",
		Category: &fhir.CodeableConcept{
			Coding: []fhir.Coding{
				o.PatientTypeCoding(),
			},
		},
		MedicationReference: fhir.Reference{
			Reference: util.StrPtrFmt("urn:uuid:%s", o.MedicationId),
			Display:   util.StrPtr(o.KfaDisplay),
		},
		Subject: &fhir.Reference{
			Reference: util.StrPtrFmt("Patient/%s", o.PatientId),
			Display:   util.StrPtr(o.PatientName),
		},
		Context: &fhir.Reference{
			Reference: util.StrPtrFmt("Encounter/%s", o.EncounterId),
		},
		WhenPrepared:   util.StrPtr(o.PreparedDate),
		WhenHandedOver: util.StrPtr(o.HandoverDate),
	}

	return medication, medicationDispense

}
func (o *MedicationDispense) BundleEntries() ([]fhir.BundleEntry, error) {
	var result []fhir.BundleEntry
	medication, medicationDispense := o.Resources()

	medicationEntry, err := BundleEntry(medication, o.MedicationId, "Medication")
	if err != nil {
		return nil, err
	}

	result = append(result, *medicationEntry)

	medicationDispenseEntry, err := BundleEntry(medicationDispense, o.MedicationDispenseId, "MedicationDispense", WithRemoveKey("medicationCodeableConcept"))
	if err != nil {
		return nil, err
	}
	result = append(result, *medicationDispenseEntry)
	return result, nil
}
