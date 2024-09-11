//go:build sleman

package simrs

import (
	"context"
	"github.com/jasoet/fhir-worker/shared/model"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	GetDiagnosisByVisitId = `
		select pv.VISIT_ID as visit_id,
			   pd.DATE_OF_DIAGNOSA as diagnosis_date,
			   pd.DIAGNOSA_ID as diagnosis_code,
			   d.NAME_OF_DIAGNOSA as diagnosis_name,
			   e.ihs_no as practitioner_satusehat_id,
			   e.FULLNAME as practitioner_name
		from PASIEN_DIAGNOSA pd
				 join PASIEN_VISITATION pv on pv.VISIT_ID = pd.VISIT_ID
				 join EMPLOYEE_ALL e on e.EMPLOYEE_ID = pd.EMPLOYEE_ID
				 join DIAGNOSA d on pd.DIAGNOSA_ID = d.DIAGNOSA_ID
		where pd.VISIT_ID = :visit_id
			`

	GetProcedureByVisitId = `
	    SELECT 1
			`

	GetVisitBetween = `
			SELECT 
				pv.VISIT_ID AS visit_id, 
				p.NO_REGISTRATION AS patient_id, 
				p.ihs_no AS patient_satusehat_id, 
				p.kip AS patient_nik,
				p.NAME_OF_PASIEN AS patient_name, 
				p.GENDER AS patient_sex,
				p.DATE_OF_BIRTH AS patient_birth_date,
				p.CONTACT_ADDRESS AS patient_address,
				pv.VISIT_DATE AS visit_date, 
				e.ihs_no AS practitioner_satusehat_id, 
				e.nik AS practitioner_nik,
				e.FULLNAME AS practitioner_name, 
				c.NAME_OF_CLINIC AS clinic_name, 
				c.id_location_satusehat AS clinic_satusehat_id, 
				rr.suhu AS temperature, 
				rr.nafas AS respiration_rate,
				rr.tensi AS blood_pressure,
				rr.nadi AS heart_rate,
				pv.VISIT_DATE AS visit_date, -- date only 
				rr.created_date AS visit_arrived_time,	
				rr.tgl_pengkajian AS visit_inprogress_date, 
				rr.jam_pengkajian AS visit_inprogress_hour, 
				rr.modi_date AS visit_end_time
			FROM 
				PASIEN_VISITATION pv
			JOIN 
				PASIEN p ON pv.NO_REGISTRATION = p.NO_REGISTRATION
			JOIN 
				CLINIC c ON pv.CLINIC_ID = c.CLINIC_ID
			JOIN 
				EMPLOYEE_ALL e ON pv.EMPLOYEE_ID = e.EMPLOYEE_ID
			JOIN 
				riwayat_rajal rr ON pv.VISIT_ID = rr.visit_id
			WHERE 
				p.ihs_no IS NOT NULL 
				AND e.ihs_no IS NOT NULL
				AND pv.VISIT_DATE between :start_date AND :end_date
			ORDER BY 
				pv.VISIT_DATE DESC;

            `

	GetObservationLabByVisitId = `
	    SELECT 1
			`

	GetObservationRadiologyByVisitId = `
	    SELECT 1
			`

	GetMedicationRequestByVisitId = `
		select bo.VISIT_ID as visit_id,
			   bo.RESEP_NO as prescription_id,
			   bo.TREAT_DATE as date,
			   bo.TREATMENT as treatment,
			   bo.aturan_pakai as usage,
			   bo.MODIFIED_DATE as modified_date,
			   bo.posting_date as posting_date,
			   g.BRAND_ID as medication_id,
			   g.NAME as medication_name,
			   e.nik as practitioner_name,
			   e.ihs_no as practitioner_satusehat_id,
			   e.FULLNAME as practitioner_name
		from bill_apotik bo
				 inner join PASIEN_VISITATION pv on bo.VISIT_ID = pv.VISIT_ID
				 inner join GOODS g on bo.BRAND_ID = g.BRAND_ID
				 inner join EMPLOYEE_ALL e on e.EMPLOYEE_ID = bo.EMPLOYEE_ID
		WHERE pv.VISIT_ID = :visit_id
    `

	GetMedicationDispenseByVisitId = `
		select bo.VISIT_ID as visit_id,
			   bo.RESEP_NO as prescription_id,
			   bo.TREAT_DATE as date,
			   bo.TREATMENT as treatment,
			   bo.aturan_pakai as usage,
			   bo.MODIFIED_DATE as modified_date,
			   bo.posting_date as posting_date,
			   g.BRAND_ID as medication_id,
			   g.NAME as medication_name,
			   e.nik as practitioner_name,
			   e.ihs_no as practitioner_satusehat_id,
			   e.FULLNAME as practitioner_name
		from bill_apotik bo
				 inner join PASIEN_VISITATION pv on bo.VISIT_ID = pv.VISIT_ID
				 inner join GOODS g on bo.BRAND_ID = g.BRAND_ID
				 inner join EMPLOYEE_ALL e on e.EMPLOYEE_ID = bo.EMPLOYEE_ID
		WHERE pv.VISIT_ID = :visit_id
    `
)

type slemanQuery struct {
	*sqlx.DB
	getVisitStmt                     *sqlx.NamedStmt
	getDiagnosisByVisitStmt          *sqlx.NamedStmt
	getMedicationRequestByVisitStmt  *sqlx.NamedStmt
	getMedicationDispenseByVisitStmt *sqlx.NamedStmt
	getProcedureByVisitStmt          *sqlx.NamedStmt
	getObservationLabByVisitId       *sqlx.NamedStmt
	getObservationRadiologyByVisitId *sqlx.NamedStmt
}

func NewQuery(pool *sqlx.DB) (Query, error) {
	queryOps := &slemanQuery{
		DB: pool,
	}

	var err error
	queryOps.getVisitStmt, err = queryOps.DB.PrepareNamed(GetVisitBetween)
	if err != nil {
		return nil, err
	}

	queryOps.getDiagnosisByVisitStmt, err = queryOps.DB.PrepareNamed(GetDiagnosisByVisitId)
	if err != nil {
		return nil, err
	}

	queryOps.getObservationLabByVisitId, err = queryOps.DB.PrepareNamed(GetObservationLabByVisitId)
	if err != nil {
		return nil, err
	}

	queryOps.getObservationRadiologyByVisitId, err = queryOps.DB.PrepareNamed(GetObservationRadiologyByVisitId)
	if err != nil {
		return nil, err
	}

	queryOps.getMedicationRequestByVisitStmt, err = queryOps.DB.PrepareNamed(GetMedicationRequestByVisitId)
	if err != nil {
		return nil, err
	}

	queryOps.getMedicationDispenseByVisitStmt, err = queryOps.DB.PrepareNamed(GetMedicationDispenseByVisitId)
	if err != nil {
		return nil, err
	}

	queryOps.getProcedureByVisitStmt, err = queryOps.DB.PrepareNamed(GetProcedureByVisitId)
	if err != nil {
		return nil, err
	}

	return queryOps, nil
}

func (f *slemanQuery) GetVisitBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.Visit, error) {
	parameter := map[string]any{
		"start_date": startDate,
		"end_date":   endDate,
	}

	var results []model.Visit

	rows, err := f.getVisitStmt.QueryxContext(ctx, parameter)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		result := make(map[string]any)
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}

		visit := BuildVisit(result)
		results = append(results, visit)
	}

	return results, nil
}

func (f *slemanQuery) GetDiagnosisByVisitId(ctx context.Context, visitId string) (model.DiagnosisList, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results model.DiagnosisList

	rows, err := f.getDiagnosisByVisitStmt.QueryxContext(ctx, parameter)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		result := make(map[string]any)
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}

		diagnosis := BuildDiagnosis(result)
		results = append(results, diagnosis)
	}

	return results, nil
}

func (f *slemanQuery) GetMedicationRequestByVisitId(ctx context.Context, visitId string) (model.MedicationRequestList, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results model.MedicationRequestList

	rows, err := f.getMedicationRequestByVisitStmt.QueryxContext(ctx, parameter)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		result := make(map[string]any)
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}

		medication := BuildMedicationRequest(result)
		results = append(results, medication)
	}

	return results, nil
}

func (f *slemanQuery) GetMedicationDispenseByVisitId(ctx context.Context, visitId string) (model.MedicationDispenseList, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []model.MedicationDispense

	err := f.getMedicationDispenseByVisitStmt.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *slemanQuery) GetProcedureByVisitId(ctx context.Context, visitId string) (model.ProcedureList, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []model.Procedure

	err := f.getProcedureByVisitStmt.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *slemanQuery) GetObservationLabByVisitId(ctx context.Context, visitId string) (model.ObservationLabList, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []model.ObservationLab

	err := f.getObservationLabByVisitId.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *slemanQuery) GetObservationRadiologyByVisitId(ctx context.Context, visitId string) (model.ObservationRadiologyList, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []model.ObservationRadiology

	err := f.getObservationLabByVisitId.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}
