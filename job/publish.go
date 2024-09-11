package job

import (
	"context"
	"errors"
	"fmt"
	"github.com/jasoet/fhir-worker/shared/model"
	"github.com/rs/zerolog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"

	"github.com/jasoet/fhir-worker/internal/db"
	"github.com/jasoet/fhir-worker/internal/entity"
	"github.com/jasoet/fhir-worker/internal/resource"
	"github.com/jasoet/fhir-worker/internal/satusehat"
	"github.com/jasoet/fhir-worker/pkg/file"
	"github.com/jasoet/fhir-worker/pkg/util"
)

type Publish struct {
	convertToUtc   bool
	simulationMode bool
	simulationDir  string
	organizationId string
	sendDelay      time.Duration
	client         *satusehat.Client
	repository     *db.Repository
}

type PublishOption func(*Publish) error

func NewPublish(options ...PublishOption) (*Publish, error) {
	p := &Publish{
		simulationDir:  os.TempDir(),
		simulationMode: false,
		convertToUtc:   false,
		sendDelay:      2 * time.Second,
	}

	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}

	if p.client == nil {
		return nil, fmt.Errorf("publish.client is required")
	}

	if p.repository == nil {
		return nil, fmt.Errorf("publish.repository is required")
	}

	return p, nil
}

func WithClientAndRepository(client *satusehat.Client, repository *db.Repository) PublishOption {
	return func(p *Publish) error {
		p.client = client
		p.repository = repository
		return nil
	}
}

func WithSimulationMode(mode bool) PublishOption {
	return func(p *Publish) error {
		p.simulationMode = mode
		return nil
	}
}

func WithSimulationDir(dir string) PublishOption {
	return func(p *Publish) error {
		p.simulationDir = dir
		return nil
	}
}

func WithConvertUtc(convert bool) PublishOption {
	return func(p *Publish) error {
		p.convertToUtc = convert
		return nil
	}
}

func WithOrganizationId(id string) PublishOption {
	return func(p *Publish) error {
		p.organizationId = id
		return nil
	}
}

func WithSendDelay(delay time.Duration) PublishOption {
	return func(p *Publish) error {
		p.sendDelay = delay
		return nil
	}
}

func (p *Publish) Process(ctx context.Context) error {
	logger := log.With().Ctx(ctx).Str("function", "Publish Process").Logger()

	internals, err := p.repository.ReadyToPublish(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch Ready data.")
		return err
	}

	logger.Info().
		Int("ready-visit-count", len(internals)).
		Bool("simulation-mode", p.simulationMode).
		Msg("publish process started")

	if p.simulationMode {
		logger = logger.With().Bool("simulation mode", p.simulationMode).Logger()
		if err := os.MkdirAll(p.simulationDir, 0755); err != nil {
			logger.Fatal().Str("simulation_dir", p.simulationDir).Err(err).Msg("failed to create simulation directory")
			return err
		}
	}

	for _, internal := range internals {
		select {
		case <-ctx.Done():
			logger.Info().Msg("context done, publish process terminated")
			return fmt.Errorf("context done, publish process terminated")
		default:
			time.Sleep(p.sendDelay)

			if err := p.processInternal(ctx, &internal, logger); err != nil {
				logger.Error().Err(err).Str("VisitId", internal.VisitID).Msg("Failed to process internal")
			}
		}
	}

	logger.Info().
		Int("ready-visit-count", len(internals)).
		Bool("simulation-mode", p.simulationMode).
		Msg("publish process finished")

	return nil
}

func (p *Publish) processInternal(ctx context.Context, internal *entity.SatuSehatInternal, logger zerolog.Logger) error {
	bundle, err := p.generateBundle(internal)
	if err != nil {
		logger.Error().Any("bundle", bundle).Any("data", internal).Err(err).Msg("generate bundle failed")
		return err
	}

	payload, err := bundle.MarshalJSON()
	if err != nil {
		logger.Error().Any("data", internal).Err(err).Msg("marshalling JSON failed")
		return err
	}

	logger.Debug().Str("VisitId", internal.VisitID).Int("payload_size", len(payload)).Msg("Processing data")

	if p.simulationMode {
		return p.simulateProcessing(internal.VisitID, payload, logger)
	}

	return p.sendToSatuSehat(ctx, internal.VisitID, payload, logger)
}

func (p *Publish) simulateProcessing(visitID string, payload []byte, logger zerolog.Logger) error {
	fileName := fmt.Sprintf("%s/%v.json", p.simulationDir, visitID)
	if err := file.WritePrettyJson(fileName, payload, 0644); err != nil {
		logger.Error().Str("file_name", fileName).Int("payload_len", len(payload)).Err(err).Msg("store payload failed")
		return err
	}

	logger.Debug().Int("payload_len", len(payload)).Str("filename", fileName).Msg("store payload data to file")
	return nil
}

func (p *Publish) sendToSatuSehat(ctx context.Context, visitID string, payload []byte, logger zerolog.Logger) error {
	logger.Debug().Int("payload_len", len(payload)).Msg("sending data to satusehat")

	requestDate := time.Now()
	respBody, err := p.client.PostBundle(ctx, string(payload))
	if err != nil {
		if errors.Is(err, &satusehat.ExecutionError{}) || errors.Is(err, &satusehat.UnauthorizedError{}) {
			logger.Error().Err(err).Msg("error when executing PostBundle")
			return err
		}

		if _, updateErr := p.repository.UpdatePublishStatus(ctx, visitID, string(payload), respBody, requestDate, entity.RequestError); updateErr != nil {
			logger.Error().Any("status", entity.RequestError).Any("payload_len", len(payload)).Err(updateErr).Msg("update database failed")
		}
		return err
	}

	if _, updateErr := p.repository.UpdatePublishStatus(ctx, visitID, string(payload), respBody, requestDate, entity.Success); updateErr != nil {
		logger.Error().Any("status", entity.Success).Any("payload_len", len(payload)).Any("response_len", len(respBody)).Err(updateErr).Msg("update database failed")
	}

	logger.Debug().Int("payload_len", len(payload)).Msg("data sent to satusehat")
	return nil
}

func (p *Publish) generateBundle(internal *entity.SatuSehatInternal) (*fhir.Bundle, error) {
	encounterUid := uuid.New().String()
	visitDetail := internal.VisitDetail()

	var entries []fhir.BundleEntry
	diagnosisEntries, encounterDiagnosis, err := p.generateDiagnosisEntries(encounterUid, visitDetail, internal.Diagnosis())
	if err != nil {
		return nil, err
	}

	encounterEntry, err := p.generateEncounterEntry(encounterUid, visitDetail, encounterDiagnosis)
	if err != nil {
		return nil, err
	}
	entries = append(entries, *encounterEntry)

	vitalSignEntries, err := p.generateVitalSignEntries(encounterUid, visitDetail, internal.VitalSign())
	if err != nil {
		return nil, err
	}

	entries = append(entries, vitalSignEntries...)

	entries = append(entries, diagnosisEntries...)

	medicationRequestEntries, err := p.generateMedicationRequestEntries(encounterUid, visitDetail, internal.MedicationRequest())
	if err != nil {
		return nil, err
	}
	entries = append(entries, medicationRequestEntries...)

	medicationDispenseEntries, err := p.generateMedicationDispenseEntries(encounterUid, visitDetail, internal.MedicationDispense())
	if err != nil {
		return nil, err
	}
	entries = append(entries, medicationDispenseEntries...)

	return &fhir.Bundle{
		Type:  fhir.BundleTypeTransaction,
		Entry: entries,
	}, err
}

func (p *Publish) generateEncounterEntry(encounterUid string, visitDetail *model.VisitDetail, encounterDiagnosis []resource.EncounterDiagnosis) (*fhir.BundleEntry, error) {
	encounter := resource.Encounter{
		EncounterId:             encounterUid,
		PatientSatuSehatId:      visitDetail.PatientSatusehatId,
		PatientName:             visitDetail.PatientName,
		OrganizationId:          p.organizationId,
		LocationName:            visitDetail.ClinicName,
		LocationId:              visitDetail.ClinicSatuSehatId,
		PeriodStartDate:         util.StdTimeToString(&visitDetail.PeriodStartDate, p.convertToUtc),
		PeriodEndDate:           util.StdTimeToString(&visitDetail.PeriodEndDate, p.convertToUtc),
		ArrivedStartTime:        util.StdTimeToString(visitDetail.ArrivedStartTime, p.convertToUtc),
		ArrivedEndTime:          util.StdTimeToString(visitDetail.ArrivedEndTime, p.convertToUtc),
		InProgressStartTime:     util.StdTimeToString(visitDetail.InProgressStartTime, p.convertToUtc),
		InProgressEndTime:       util.StdTimeToString(visitDetail.InProgressEndTime, p.convertToUtc),
		FinishStartTime:         util.StdTimeToString(visitDetail.FinishStartTime, p.convertToUtc),
		FinishEndTime:           util.StdTimeToString(visitDetail.FinishEndTime, p.convertToUtc),
		PractitionerSatuSehatId: visitDetail.PractitionerId,
		PractitionerName:        visitDetail.PractitionerName,
		Diagnosis:               encounterDiagnosis,
	}
	return encounter.BundleEntry()
}

func (p *Publish) generateVitalSignEntries(encounterUid string, visitDetail *model.VisitDetail, vitalSign *model.VitalSign) ([]fhir.BundleEntry, error) {
	vitalSignResources := &resource.VitalSign{
		EncounterId:             encounterUid,
		SystoleId:               uuid.New().String(),
		DiastoleId:              uuid.New().String(),
		HeartRateId:             uuid.New().String(),
		TemperatureId:           uuid.New().String(),
		RespirationRateId:       uuid.New().String(),
		OxygenSaturationId:      uuid.New().String(),
		PatientSatuSehatId:      visitDetail.PatientSatusehatId,
		PatientName:             visitDetail.PatientName,
		Time:                    util.StdTimeToString(&visitDetail.PeriodStartDate, p.convertToUtc),
		PractitionerSatuSehatId: visitDetail.PractitionerId,
		PractitionerName:        visitDetail.PractitionerName,
		Systole:                 vitalSign.Systole,
		Diastole:                vitalSign.Diastole,
		HeartRate:               vitalSign.HeartRate,
		Temperature:             vitalSign.Temperature,
		RespirationRate:         vitalSign.RespirationRate,
		OxygenSaturation:        vitalSign.OxygenSaturation,
	}
	return vitalSignResources.BundleEntries()
}

func (p *Publish) generateDiagnosisEntries(encounterUid string, visitDetail *model.VisitDetail, diagnosisList *model.DiagnosisList) ([]fhir.BundleEntry, []resource.EncounterDiagnosis, error) {
	var encounterDiagnosis []resource.EncounterDiagnosis
	var entries []fhir.BundleEntry
	if diagnosisList != nil {
		for _, diagnosis := range *diagnosisList {
			if !diagnosis.Invalid() {
				conditionId := uuid.New().String()
				conditionDisplay := diagnosis.DiagnosisName

				conditionDiagnosis := resource.ConditionDiagnosis{
					ConditionId:        conditionId,
					EncounterId:        encounterUid,
					PatientSatuSehatId: visitDetail.PatientSatusehatId,
					PatientName:        visitDetail.PatientName,
					Time:               util.StdTimeToString(&diagnosis.DiagnosisDate, p.convertToUtc),
					IcdCode:            diagnosis.DiagnosisCode,
					IcdName:            diagnosis.DiagnosisName,
				}

				diagnosisEntry, err := conditionDiagnosis.BundleEntry()
				if err != nil {
					return nil, nil, err
				}
				encounterDiagnosis = append(encounterDiagnosis, resource.EncounterDiagnosis{
					Id:      conditionId,
					Display: conditionDisplay,
				})
				entries = append(entries, *diagnosisEntry)
			}
		}
	}
	return entries, encounterDiagnosis, nil
}

func (p *Publish) generateMedicationRequestEntries(encounterUid string, visitDetail *model.VisitDetail, medicationRequestList *model.MedicationRequestList) ([]fhir.BundleEntry, error) {
	var entries []fhir.BundleEntry
	if medicationRequestList != nil {
		for _, request := range *medicationRequestList {
			if !request.Invalid() {
				res := resource.MedicationRequest{
					MedicationId:        uuid.New().String(),
					MedicationRequestId: uuid.New().String(),
					EncounterId:         encounterUid,
					PatientId:           visitDetail.PatientSatusehatId,
					PatientName:         visitDetail.PatientName,
					OrganizationId:      p.organizationId,
					Date:                util.StdTimeToString(request.Date, p.convertToUtc),
					PractitionerId:      util.StringNotNil(request.PractitionerId),
					PractitionerName:    util.StringNotNil(request.PractitionerName),
					PrescriptionId:      util.IntToString(&request.PrescriptionId),
					KfaCode:             util.StringNotNil(request.KfaCode),
					KfaDisplay:          util.StringNotNil(request.KfaName),
					Type:                request.Type,
					PatientType:         request.PatientType,
				}

				requestList, err := res.BundleEntries()
				if err != nil {
					return nil, err
				}
				entries = append(entries, requestList...)
			}
		}
	}
	return entries, nil
}

func (p *Publish) generateMedicationDispenseEntries(encounterUid string, visitDetail *model.VisitDetail, medicationDispenseList *model.MedicationDispenseList) ([]fhir.BundleEntry, error) {
	var entries []fhir.BundleEntry
	if medicationDispenseList != nil {
		for _, dispense := range *medicationDispenseList {
			if !dispense.Invalid() {
				res := resource.MedicationDispense{
					MedicationId:         uuid.New().String(),
					MedicationDispenseId: uuid.New().String(),
					EncounterId:          encounterUid,
					PatientId:            visitDetail.PatientSatusehatId,
					PatientName:          visitDetail.PatientName,
					OrganizationId:       p.organizationId,
					PractitionerId:       util.StringNotNil(dispense.PractitionerId),
					PractitionerName:     util.StringNotNil(dispense.PractitionerName),
					PrescriptionId:       util.IntToString(&dispense.PrescriptionId),
					KfaCode:              util.StringNotNil(dispense.KfaCode),
					KfaDisplay:           util.StringNotNil(dispense.KfaName),
					Type:                 dispense.Type,
					PatientType:          dispense.PatientType,
					PreparedDate:         util.StdTimeToString(dispense.PrescriptionStartDate, p.convertToUtc),
					HandoverDate:         util.StdTimeToString(dispense.HandoverDate, p.convertToUtc),
					BatchNumber:          dispense.BatchNumber,
					ExpirationDate:       util.StdTimeToString(dispense.ExpiredDate, p.convertToUtc),
				}

				dispenseList, err := res.BundleEntries()
				if err != nil {
					return nil, err
				}
				entries = append(entries, dispenseList...)
			}
		}
	}
	return entries, nil
}
