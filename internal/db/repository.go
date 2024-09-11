package db

import (
	"context"
	"database/sql"
	"github.com/jasoet/fhir-worker/internal/entity"
	"github.com/jasoet/fhir-worker/pkg/util"
	shared "github.com/jasoet/fhir-worker/shared/model"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

const (
	GetByStatus = `
		SELECT 
			si.visit_id, 
			si.visit_date,
			si.satusehat_patient_id, 
			si.visit_detail, 
			si.vital_sign, 
			si.diagnosis, 
			si.lab, 
			si.radiology, 
			si.medication_request, 
			si.medication_dispense, 
			si.medical_procedure, 
			si.publish_date, 
			si.publish_request, 
			si.publish_response, 
			si.publish_status, 
			si.mapping_errors,
			si.mapping_status 
		FROM 
			satusehat AS si
		WHERE
			si.mapping_status = :mapping_status;
   `

	Insert = `
		INSERT INTO satusehat (
			visit_id, 
		    visit_date,
			satusehat_patient_id, 
			visit_detail, 
			vital_sign, 
			publish_status, 
			mapping_status,
		    mapping_errors
		) 
		VALUES (
			:visit_id, 
		    :visit_date,
			:satusehat_patient_id, 
			:visit_detail, 
			:vital_sign, 
			:publish_status, 
			:mapping_status,
		    :mapping_errors
		);
   `

	UpdateDiagnosis = `
		UPDATE satusehat
		SET diagnosis = :diagnosis
		WHERE visit_id = :visit_id;
	`

	UpdateLab = `
		UPDATE satusehat
		SET lab = :lab
		WHERE visit_id = :visit_id;
	`

	UpdateRadiology = `
		UPDATE satusehat
		SET radiology = :radiology
		WHERE visit_id = :visit_id;
	`

	UpdateMedicationRequest = `
		UPDATE satusehat
		SET medication_request = :medication_request
		WHERE visit_id = :visit_id;
	`

	UpdateMedicationDispense = `
		UPDATE satusehat
		SET medication_dispense = :medication_dispense
		WHERE visit_id = :visit_id;
	`

	UpdateMedicalProcedure = `
		UPDATE satusehat
		SET medical_procedure = :medical_procedure
		WHERE visit_id = :visit_id;
	`

	UpdatePublishStatus = `
		UPDATE satusehat
		SET publish_response = :publish_response,
			publish_request = :publish_request,
		    publish_date = :publish_date,
			publish_status = :publish_status
		WHERE visit_id = :visit_id;
	`

	UpdateMappingStatus = `
		UPDATE satusehat
		SET mapping_status = :mapping_status
		WHERE visit_id = :visit_id;
	`

	UpdateMappingErrors = `
		UPDATE satusehat
		SET mapping_errors = :mapping_errors
		WHERE visit_id = :visit_id;
	`

	IsExists = `
        SELECT count(visit_id) FROM satusehat WHERE visit_id = :visit_id;
	`
)

type Repository struct {
	db                       *sqlx.DB
	insert                   *sqlx.NamedStmt
	isExists                 *sqlx.NamedStmt
	getByStatus              *sqlx.NamedStmt
	updateDiagnosis          *sqlx.NamedStmt
	updateLab                *sqlx.NamedStmt
	updateRadiology          *sqlx.NamedStmt
	updateMedicationRequest  *sqlx.NamedStmt
	updateMedicationDispense *sqlx.NamedStmt
	updateMedicalProcedure   *sqlx.NamedStmt
	updatePublishStatus      *sqlx.NamedStmt
	updateMappingStatus      *sqlx.NamedStmt
	updateMappingErrors      *sqlx.NamedStmt
	mu                       sync.Mutex // Mutex for thread-safety
}

func newRepository(db *sqlx.DB) (*Repository, error) {
	insertNewStmt, err := db.PrepareNamed(Insert)
	if err != nil {
		return nil, err
	}

	isExists, err := db.PrepareNamed(IsExists)
	if err != nil {
		return nil, err
	}

	getByStatusStmt, err := db.PrepareNamed(GetByStatus)
	if err != nil {
		return nil, err
	}

	updateDiagnosisStmt, err := db.PrepareNamed(UpdateDiagnosis)
	if err != nil {
		return nil, err
	}

	updateLabStmt, err := db.PrepareNamed(UpdateLab)
	if err != nil {
		return nil, err
	}

	updateRadiologyStmt, err := db.PrepareNamed(UpdateRadiology)
	if err != nil {
		return nil, err
	}

	updateMedicationRequestStmt, err := db.PrepareNamed(UpdateMedicationRequest)
	if err != nil {
		return nil, err
	}

	updateMedicationDispenseStmt, err := db.PrepareNamed(UpdateMedicationDispense)
	if err != nil {
		return nil, err
	}

	updateMedicalProcedureStmt, err := db.PrepareNamed(UpdateMedicalProcedure)
	if err != nil {
		return nil, err
	}

	updatePublishStatusStmt, err := db.PrepareNamed(UpdatePublishStatus)
	if err != nil {
		return nil, err
	}

	updateMappingStatusStmt, err := db.PrepareNamed(UpdateMappingStatus)
	if err != nil {
		return nil, err
	}

	updateMappingErrorsStmt, err := db.PrepareNamed(UpdateMappingErrors)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:                       db,
		insert:                   insertNewStmt,
		isExists:                 isExists,
		getByStatus:              getByStatusStmt,
		updateDiagnosis:          updateDiagnosisStmt,
		updateLab:                updateLabStmt,
		updateRadiology:          updateRadiologyStmt,
		updateMedicationRequest:  updateMedicationRequestStmt,
		updateMedicationDispense: updateMedicationDispenseStmt,
		updateMedicalProcedure:   updateMedicalProcedureStmt,
		updatePublishStatus:      updatePublishStatusStmt,
		updateMappingStatus:      updateMappingStatusStmt,
		updateMappingErrors:      updateMappingErrorsStmt,
		mu:                       sync.Mutex{},
	}, nil
}

func (r *Repository) IsExists(ctx context.Context, visitId string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int

	err := r.isExists.GetContext(ctx, &count, map[string]any{
		"visit_id": visitId,
	})

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) ReadyToPublish(ctx context.Context) ([]entity.SatuSehatInternal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	parameter := map[string]any{
		"mapping_status": entity.Ready,
	}

	var results []entity.SatuSehatInternal

	err := r.getByStatus.SelectContext(ctx, &results, parameter)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) Incomplete(ctx context.Context) ([]entity.SatuSehatInternal, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	parameter := map[string]any{
		"mapping_status": entity.Incomplete,
	}

	var results []entity.SatuSehatInternal

	err := r.getByStatus.SelectContext(ctx, &results, parameter)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) InsertValid(
	ctx context.Context,
	visitId string,
	visitDate time.Time,
	satusehatPatientId string,
	visitDetail shared.VisitDetail,
	visitSign shared.VitalSign,
) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.insert.ExecContext(ctx, map[string]any{
		"visit_id":             visitId,
		"visit_date":           visitDate,
		"satusehat_patient_id": satusehatPatientId,
		"visit_detail":         util.MarshalToJson(visitDetail),
		"vital_sign":           util.MarshalToJson(visitSign),
		"publish_status":       entity.Preparing,
		"mapping_status":       entity.Incomplete,
		"mapping_errors":       "",
	})
}

func (r *Repository) InsertInvalid(
	ctx context.Context,
	visitId string,
	visitDate time.Time,
	satusehatPatientId string,
	visitDetail shared.VisitDetail,
	visitSign shared.VitalSign,
	mappingErrors string,
) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.insert.ExecContext(ctx, map[string]any{
		"visit_id":             visitId,
		"visit_date":           visitDate,
		"satusehat_patient_id": satusehatPatientId,
		"visit_detail":         util.MarshalToJson(visitDetail),
		"vital_sign":           util.MarshalToJson(visitSign),
		"mapping_errors":       mappingErrors,
		"publish_status":       entity.Preparing,
		"mapping_status":       entity.Invalid,
	})
}

func (r *Repository) UpdateDiagnosis(ctx context.Context, visitId string, diagnosis []shared.Diagnosis) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateDiagnosis.ExecContext(ctx, map[string]any{
		"visit_id":  visitId,
		"diagnosis": util.MarshalToJson(diagnosis),
	})
}

func (r *Repository) UpdateLab(ctx context.Context, visitId string, lab []shared.ObservationLab) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateLab.ExecContext(ctx, map[string]any{
		"visit_id": visitId,
		"lab":      util.MarshalToJson(lab),
	})
}

func (r *Repository) UpdateRadiology(ctx context.Context, visitId string, radiology []shared.ObservationRadiology) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateRadiology.ExecContext(ctx, map[string]any{
		"visit_id":  visitId,
		"radiology": util.MarshalToJson(radiology),
	})
}

func (r *Repository) UpdateMedicationRequest(ctx context.Context, visitId string, medicationRequest []shared.MedicationRequest) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateMedicationRequest.ExecContext(ctx, map[string]any{
		"visit_id":           visitId,
		"medication_request": util.MarshalToJson(medicationRequest),
	})
}

func (r *Repository) UpdateMedicationDispense(ctx context.Context, visitId string, medicationDispense []shared.MedicationDispense) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateMedicationDispense.ExecContext(ctx, map[string]any{
		"visit_id":            visitId,
		"medication_dispense": util.MarshalToJson(medicationDispense),
	})
}

func (r *Repository) UpdateMedicalProcedure(ctx context.Context, visitId string, medicalProcedure []shared.Procedure) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateMedicalProcedure.ExecContext(ctx, map[string]any{
		"visit_id":          visitId,
		"medical_procedure": util.MarshalToJson(medicalProcedure),
	})
}

func (r *Repository) UpdatePublishStatus(ctx context.Context, visitId string, publishRequest string, publishResponse string, publishDate time.Time, status entity.PublishStatus) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updatePublishStatus.ExecContext(ctx, map[string]any{
		"visit_id":         visitId,
		"publish_response": publishResponse,
		"publish_date":     publishDate,
		"publish_status":   status,
		"publish_request":  publishRequest,
	})
}

func (r *Repository) UpdateMappingStatus(ctx context.Context, visitId string, mappingStatus entity.MappingStatus) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateMappingStatus.ExecContext(ctx, map[string]any{
		"visit_id":       visitId,
		"mapping_status": mappingStatus,
	})

}
func (r *Repository) UpdateMappingErrors(ctx context.Context, visitId string, mappingErrors string) (sql.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.updateMappingErrors.ExecContext(ctx, map[string]any{
		"visit_id":       visitId,
		"mapping_errors": mappingErrors,
	})

}
