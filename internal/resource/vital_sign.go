package resource

import (
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type VitalSign struct {
	EncounterId             string `validate:"required"`
	SystoleId               string `validate:"required"`
	DiastoleId              string `validate:"required"`
	HeartRateId             string `validate:"required"`
	TemperatureId           string `validate:"required"`
	OxygenSaturationId      string `validate:"required"`
	RespirationRateId       string `validate:"required"`
	PatientSatuSehatId      string `validate:"required"`
	PatientName             string `validate:"required"`
	Time                    string `validate:"required"`
	PractitionerSatuSehatId string
	PractitionerName        string
	Systole                 string
	Diastole                string
	HeartRate               string
	Temperature             string
	RespirationRate         string
	OxygenSaturation        string
}

func (o *VitalSign) BundleEntries() ([]fhir.BundleEntry, error) {
	var bundleEntries []fhir.BundleEntry
	observations := o.Observations()
	for _, obs := range observations {
		bundleEntry, err := obs.BundleEntry()
		if err != nil {
			return nil, err
		}
		bundleEntries = append(bundleEntries, *bundleEntry)
	}
	return bundleEntries, nil
}

func (o *VitalSign) Observations() []Observation {
	var observations []Observation
	if util.StringNotEmpty(o.Systole) {
		obsrv := Observation{
			EncounterId:             o.EncounterId,
			ObservationId:           o.SystoleId,
			PatientSatuSehatId:      o.PatientSatuSehatId,
			PatientName:             o.PatientName,
			Time:                    o.Time,
			PractitionerName:        o.PractitionerName,
			PractitionerSatuSehatId: o.PractitionerSatuSehatId,
			LoincCode:               "8480-6",
			LoincDisplay:            "Systolic blood pressure",
			ValueQuantity: &ObservationValue{
				Value: o.Systole,
				Unit:  "mmHg",
				Code:  "mm[Hg]",
			},
		}

		observations = append(observations, obsrv)

	}

	if util.StringNotEmpty(o.Diastole) {
		obsrv := Observation{
			ObservationId:           o.DiastoleId,
			EncounterId:             o.EncounterId,
			PatientSatuSehatId:      o.PatientSatuSehatId,
			PatientName:             o.PatientName,
			PractitionerName:        o.PractitionerName,
			PractitionerSatuSehatId: o.PractitionerSatuSehatId,
			Time:                    o.Time,
			LoincCode:               "8462-4",
			LoincDisplay:            "Diastolic blood pressure",
			ValueQuantity: &ObservationValue{
				Value: o.Diastole,
				Unit:  "mmHg",
				Code:  "mm[Hg]",
			},
		}

		observations = append(observations, obsrv)
	}

	if util.StringNotEmpty(o.Temperature) {
		obsrv := Observation{
			ObservationId:           o.TemperatureId,
			EncounterId:             o.EncounterId,
			PatientSatuSehatId:      o.PatientSatuSehatId,
			PatientName:             o.PatientName,
			PractitionerName:        o.PractitionerName,
			PractitionerSatuSehatId: o.PractitionerSatuSehatId,
			Time:                    o.Time,
			LoincCode:               "8310-5",
			LoincDisplay:            "Body temperature",
			ValueQuantity: &ObservationValue{
				Value: o.Temperature,
				Unit:  "C",
				Code:  "Cel",
			},
		}
		observations = append(observations, obsrv)

	}

	if util.StringNotEmpty(o.HeartRate) {
		obsrv := Observation{
			ObservationId:           o.HeartRateId,
			EncounterId:             o.EncounterId,
			PatientSatuSehatId:      o.PatientSatuSehatId,
			PatientName:             o.PatientName,
			PractitionerName:        o.PractitionerName,
			PractitionerSatuSehatId: o.PractitionerSatuSehatId,
			Time:                    o.Time,
			LoincCode:               "8867-4",
			LoincDisplay:            "Heart rate",
			ValueQuantity: &ObservationValue{
				Value: o.HeartRate,
				Unit:  "beats/min",
				Code:  "/min",
			},
		}
		observations = append(observations, obsrv)

	}

	if util.StringNotEmpty(o.RespirationRate) {
		obsrv := Observation{
			ObservationId:           o.RespirationRateId,
			EncounterId:             o.EncounterId,
			PatientSatuSehatId:      o.PatientSatuSehatId,
			PatientName:             o.PatientName,
			PractitionerName:        o.PractitionerName,
			PractitionerSatuSehatId: o.PractitionerSatuSehatId,
			Time:                    o.Time,
			LoincCode:               "9279-1",
			LoincDisplay:            "Respiratory rate",
			ValueQuantity: &ObservationValue{
				Value: o.RespirationRate,
				Unit:  "breaths/min",
				Code:  "/min",
			},
		}
		observations = append(observations, obsrv)

	}

	if util.StringNotEmpty(o.OxygenSaturation) {
		obsrv := Observation{
			ObservationId:           o.OxygenSaturationId,
			EncounterId:             o.EncounterId,
			PatientSatuSehatId:      o.PatientSatuSehatId,
			PatientName:             o.PatientName,
			PractitionerName:        o.PractitionerName,
			PractitionerSatuSehatId: o.PractitionerSatuSehatId,
			Time:                    o.Time,
			LoincCode:               "2708-6",
			LoincDisplay:            "Oxygen saturation in Arterial blood",
			ValueQuantity: &ObservationValue{
				Value: o.OxygenSaturation,
				Unit:  "%",
				Code:  "%",
			},
		}
		observations = append(observations, obsrv)

	}

	return observations
}
