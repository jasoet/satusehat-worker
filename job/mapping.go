package job

import (
	"context"
	"fmt"
	"github.com/jasoet/fhir-worker/internal/db"
	"github.com/jasoet/fhir-worker/internal/entity"
	"github.com/jasoet/fhir-worker/simrs"
	"github.com/rs/zerolog/log"
	"time"
)

type Mapping struct {
	markCompleteDays  int
	lastVisitDays     int
	DisableDiagnosis  bool
	DisableLab        bool
	DisableRadiology  bool
	DisableProcedure  bool
	DisableMedication bool
	queryOps          simrs.Query
	repository        *db.Repository
}

type MappingOption func(o *Mapping) error

func WithDisableConfigs(
	disableDiagnosis bool,
	disableLab bool,
	disableRadiology bool,
	disableProcedure bool,
	disableMedication bool,
) MappingOption {
	return func(o *Mapping) error {
		o.DisableDiagnosis = disableDiagnosis
		o.DisableLab = disableLab
		o.DisableRadiology = disableRadiology
		o.DisableProcedure = disableProcedure
		o.DisableMedication = disableMedication
		return nil
	}
}

func WithConfigDays(markCompleteDays int, lastVisitDays int) MappingOption {
	return func(o *Mapping) error {
		o.markCompleteDays = markCompleteDays
		o.lastVisitDays = lastVisitDays
		return nil
	}

}

func WithQueryAndRepository(queryOps simrs.Query, repository *db.Repository) MappingOption {
	return func(o *Mapping) error {
		o.queryOps = queryOps
		o.repository = repository
		return nil
	}
}

func NewMapping(options ...MappingOption) (*Mapping, error) {
	mapping := &Mapping{
		markCompleteDays: 7,
		lastVisitDays:    7,
	}

	var err error
	for _, option := range options {
		err = option(mapping)
		if err != nil {
			return nil, err
		}
	}

	if mapping.queryOps == nil {
		return nil, fmt.Errorf("mapping.queryOps is required")
	}

	if mapping.repository == nil {
		return nil, fmt.Errorf("mapping.repository is required")
	}

	return mapping, nil
}

func (j *Mapping) CheckComplete(ctx context.Context) error {
	_log := log.With().Ctx(ctx).Str("function", "CheckComplete").
		Logger()

	internals, err := j.repository.Incomplete(ctx)
	if err != nil {
		_log.Error().Err(err).
			Msg("Failed to fetch incomplete data.")
		return err
	}

	for _, internal := range internals {
		visitDate := internal.VisitDate
		complete := false
		if time.Since(visitDate) > time.Duration(j.markCompleteDays)*24*time.Hour {
			complete = true
		}

		if !internal.Diagnosis().Invalid() && !internal.Lab().Invalid() && !internal.Radiology().Invalid() &&
			!internal.MedicationRequest().Invalid() && !internal.MedicationDispense().Invalid() && !internal.Procedure().Invalid() {
			complete = true
		}

		if complete {
			_, err := j.repository.UpdateMappingStatus(ctx, internal.VisitID, entity.Ready)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", internal.VisitID).
					Msg("Failed to update mapping status to 'Ready.'")
				continue
			}

			_log.Debug().Str("visit-id", internal.VisitID).
				Msg("Successfully updated mapping status to 'Ready.'")
		}
	}

	return nil
}

func (j *Mapping) FillVisit(ctx context.Context) error {
	_log := log.With().Ctx(ctx).Str("function", "FillVisit").
		Logger()

	internals, err := j.repository.Incomplete(ctx)
	if err != nil {
		_log.Error().Err(err).
			Msg("Failed to fetch incomplete data.")
		return err
	}

	_log.Info().
		Int("visit-count", len(internals)).
		Bool("diagnosis-disabled", j.DisableDiagnosis).
		Bool("lab-disabled", j.DisableLab).
		Bool("radiology-disabled", j.DisableRadiology).
		Bool("procedure-disabled", j.DisableProcedure).
		Bool("medication-disabled", j.DisableMedication).
		Msg("fill visit data job started")

	for _, internal := range internals {
		visitId := internal.VisitID

		_log = _log.With().
			Str("visit-id", visitId).Logger()

		if !j.DisableDiagnosis && internal.Diagnosis().Invalid() {
			_log.Debug().
				Msg("process diagnosis data")

			diagnosis, err := j.queryOps.GetDiagnosisByVisitId(ctx, visitId)
			if err != nil {
				_log.Error().Err(err).
					Msg("Failed to fetch diagnosis data.")
				continue
			}

			_, err = j.repository.UpdateDiagnosis(ctx, visitId, diagnosis)
			if err != nil {
				_log.Error().Err(err).
					Msg("Failed to update diagnosis data.")
				continue
			}

			_log.Debug().Str("visit-id", visitId).
				Msg("Successfully updated diagnosis data.")
		}

		if !j.DisableLab && internal.Lab().Invalid() {
			labs, err := j.queryOps.GetObservationLabByVisitId(ctx, visitId)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to fetch labs data.")
				continue
			}

			_, err = j.repository.UpdateLab(ctx, visitId, labs)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to update labs data.")
				continue
			}

			_log.Debug().Str("visit-id", visitId).
				Msg("Successfully updated labs data.")
		}

		if !j.DisableRadiology && internal.Radiology().Invalid() {
			radiologyData, err := j.queryOps.GetObservationRadiologyByVisitId(ctx, visitId)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to fetch radiology data.")
				continue
			}

			_, err = j.repository.UpdateRadiology(ctx, visitId, radiologyData)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to update radiology data.")
				continue
			}

			_log.Debug().Str("visit-id", visitId).
				Msg("Successfully updated radiology data.")
		}

		if !j.DisableMedication && internal.MedicationRequest().Invalid() {
			data, err := j.queryOps.GetMedicationRequestByVisitId(ctx, visitId)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to fetch MedicationRequest data.")
				continue
			}

			_, err = j.repository.UpdateMedicationRequest(ctx, visitId, data)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to update MedicationRequest data.")
				continue
			}

			_log.Debug().Str("visit-id", visitId).
				Msg("Successfully updated MedicationRequest data.")
		}

		if !j.DisableMedication && internal.MedicationDispense().Invalid() {
			data, err := j.queryOps.GetMedicationDispenseByVisitId(ctx, visitId)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to fetch MedicationDispense data.")
				continue
			}

			_, err = j.repository.UpdateMedicationDispense(ctx, visitId, data)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to update MedicationDispense data.")
				continue
			}

			_log.Debug().Str("visit-id", visitId).
				Msg("Successfully updated MedicationDispense data.")
		}

		if !j.DisableProcedure && internal.Procedure().Invalid() {
			data, err := j.queryOps.GetProcedureByVisitId(ctx, visitId)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to fetch procedure data.")
				continue
			}

			_, err = j.repository.UpdateMedicalProcedure(ctx, visitId, data)
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to update procedure data.")
				continue
			}

			_log.Debug().Str("visit-id", visitId).
				Msg("Successfully updated procedure data.")
		}

	}

	_log.Info().
		Int("visit-count", len(internals)).
		Bool("diagnosis-disabled", j.DisableDiagnosis).
		Bool("lab-disabled", j.DisableLab).
		Bool("radiology-disabled", j.DisableRadiology).
		Bool("procedure-disabled", j.DisableProcedure).
		Bool("medication-disabled", j.DisableMedication).
		Msg("fill visit data job finished")

	return nil
}

func (j *Mapping) FetchVisit(ctx context.Context) error {
	startTime := time.Now().AddDate(0, 0, -j.lastVisitDays)
	endTime := time.Now().AddDate(0, 0, 1)

	_log := log.With().Ctx(ctx).Str("function", "FetchVisit").
		Time("startTime", startTime).Time("endTime", endTime).
		Logger()

	visits, err := j.queryOps.GetVisitBetween(ctx, startTime, endTime)
	if err != nil {
		_log.Error().Err(err).
			Msg("Failed to fetch visits.")
		return err
	}

	_log.Info().
		Int("visit-count", len(visits)).
		Msg("fetch visit data job started")

	for _, visit := range visits {
		visitId := visit.VisitID
		satusehatId := visit.PatientSatusehatID

		exists, err := j.repository.IsExists(ctx, visitId)
		if err != nil {
			_log.Error().Err(err).Str("visit-id", visitId).
				Msg("Visit check failed.")
			continue
		}

		if exists {
			_log.Debug().Str("visit-id", visitId).
				Msg("Visit exists, skipping...")

			continue
		}

		validationErrors := visit.VisitDetail().Invalid()

		if validationErrors != nil {
			_log.Debug().Str("visit-id", visitId).
				Any("VisitDetail", visit.VisitDetail()).
				Msg("Visit is invalid.")

			_, err := j.repository.InsertInvalid(ctx, visitId, visit.PeriodStartDate, satusehatId, visit.VisitDetail(), visit.VitalSign(), validationErrors.Error())
			if err != nil {
				_log.Error().Err(err).Str("visit-id", visitId).
					Msg("Failed to save invalid visit data.")
				continue
			}

			_log.Debug().Str("visit-id", visitId).
				Msg("Saved visit successfully, but with 'Invalid' status.")
			continue
		}

		_, err = j.repository.InsertValid(ctx, visitId, visit.PeriodStartDate, satusehatId, visit.VisitDetail(), visit.VitalSign())
		if err != nil {
			_log.Error().Err(err).Str("visit-id", visitId).
				Msg("Failed to save visit data.")
			continue
		}

		_log.Debug().Str("visit-id", visitId).
			Msg("Successfully saved visit data.")
	}

	_log.Info().
		Int("visit-count", len(visits)).
		Msg("fetch visit data job finished")

	return nil
}
