//go:build !sleman

package simrs

import (
	"context"
	"github.com/jasoet/fhir-worker/shared/model"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	GetDiagnosisByVisitId = `
			SELECT v.id AS visit_id,
				   ad.icd_code AS diagnosis_code,
				   ad.icd_name AS diagnosis_name,
				   ad.case AS diagnosis_case,
				   ad.jenis AS diagnosis_type,
				   ad.date as diagnosis_date
			FROM visits v
					 JOIN anamnese_diagnoses ad ON (ad.visit_id = v.id)
			WHERE v.id = :visit_id
			ORDER BY ad.ordering
			`

	GetProcedureByVisitId = `
			SELECT
				v.id as visit_id,
				vicd.icd9cm_code AS procedure_code,
				vicd.icd9cm_name AS procedure_name
			FROM 
				visits v
				JOIN visits_icd9cm vicd ON (vicd.visit_id = v.id)
			WHERE v.id=:visit_id
			`

	GetVisitBetween = `
            SELECT 
                v.id AS visit_id,
                p.satusehat_patient_id AS patient_satusehat_id,
                p.nik AS patient_nik,
                p.name AS patient_name,
                p.sex AS patient_sex,
                p.birth_date AS patient_birth_date,
                p.address AS address,
                rv.name AS desa,
                rsd.name AS kecamatan,
                rd.name AS kota,
                re.name AS education_name,
                rj.name AS job_name,
                rmt.name AS marital_status_name,
                p.phone AS patient_phone,
                rpar.nik AS practitioner_nik,
                rpar.satusehat_practitioner_id AS practitioner_satusehat_id,
                rpar.name AS practitioner_name,
                rc.satusehat_location_id AS clinic_location_id,
                rc.name AS clinic_name,
                v.sistole AS visit_sistole,
                v.diastole AS visit_diastole,
                v.heart_rate AS visit_heart_rate,
                v.respiration_rate AS visit_respiration_rate,
                v.temperature AS visit_temperature,
                v.spo2 AS visit_spo2,
                v.date AS visit_date,
                v.registration_start_date AS registration_start_date,
                v.registration_date AS registration_date,
                v.pemeriksaan_start_date AS examination_start_date,
                v.pemeriksaan_end_date AS examination_end_date
			FROM 
                visits v
                JOIN patients p ON (p.id = v.patient_id )
                JOIN ref_clinics rc ON (rc.id = v.clinic_id)
                JOIN ref_paramedics rpar ON (rpar.id = v.paramedic_id)

                LEFT JOIN ref_villages rv ON (rv.id = p.village_id)
                LEFT JOIN ref_sub_districts rsd ON (rsd.id = rv.sub_district_id)
                LEFT JOIN ref_districts rd ON (rd.id = rsd.district_id)

                LEFT JOIN ref_educations re ON (re.id = p.education_id)
                LEFT JOIN ref_jobs rj ON (rj.id = p.job_id)
                LEFT JOIN ref_marital_status rmt ON (rmt.id = p.marital_status_id)
 			WHERE
				 v.date BETWEEN :start_date AND :end_date 
--                 AND rpar.satusehat_practitioner_id IS NOT NULL
--                 AND rc.satusehat_location_id IS NOT NULL
                AND (v.continue_id IS NULL OR v.continue_id<>10)
                AND v.admission_type_id<>4
                AND v.clinic_id<>1
                AND rc.type='rawat jalan'
			GROUP BY v.id
            `

	GetObservationLabByVisitId = `
			SELECT kd.visit_id AS visit_id,
				   kd.name AS lab_name,
				   kd.lab_parameter AS lab_parameter,
				   kd.lab_satuan AS lab_unit,
				   kd.lab_normal AS lab_normal,
				   kd.lab_result AS lab_result,
				   kd.lab_flag AS lab_flag,
				   kd.lab_metode AS lab_method,
				   rkbi.lab_loinc_code AS lab_loinc_code,
				   rkbi.lab_loinc_name AS lab_loinc_name,
				   rp.satusehat_practitioner_id AS practitioner_id,
				   rp.name AS practitioner_name
			FROM visits v
					 JOIN kwitansi_detail kd ON (kd.visit_id = v.id)
					 JOIN ref_komponen_biaya_item rkbi ON (rkbi.id = kd.komponen_biaya_item_id)
					 JOIN ref_paramedics rp ON (rp.id = v.paramedic_id)
			WHERE v.parent_id=:visit_id AND rkbi.lab_group_id IS NOT NULL 
			`

	GetObservationRadiologyByVisitId = `
			SELECT kd.visit_id AS visit_id,
				   kd.name AS lab_name,
				   kd.lab_parameter AS lab_parameter,
				   kd.lab_satuan AS lab_unit,
				   kd.lab_normal AS lab_normal,
				   kd.lab_result AS lab_result,
				   kd.lab_flag AS lab_flag,
				   kd.lab_metode AS lab_method,
				   rkbi.lab_loinc_code AS lab_loinc_code,
				   rkbi.lab_loinc_name AS lab_loinc_name,
				   rp.satusehat_practitioner_id AS practitioner_id,
				   rp.name AS practitioner_name
			FROM visits v
					 JOIN kwitansi_detail kd ON (kd.visit_id = v.id)
					 JOIN ref_komponen_biaya_item rkbi ON (rkbi.id = kd.komponen_biaya_item_id)
					 JOIN ref_paramedics rp ON (rp.id = v.paramedic_id)
			WHERE v.parent_id=:visit_id AND rkbi.radiologi_group_id IS NOT NULL
			`

	GetMedicationRequestByVisitId = `
			SELECT
				pt.visit_id as visit_id,
				pt.jenis_pasien as patient_type,
				pt.date as date,
				rd.code as drug_code,
				pt.id as prescription_id,
				ptd.id as prescription_detail_id,
				rd.satusehat_kfa_code as kfa_code,
				rd.satusehat_kfa_name as kfa_name,
				ptd.jenis as type,
				rp.satusehat_practitioner_id as practitioner_id,
				rp.name as paramedic_name,
				ptd.jumlah as amount,
				ptd.satuan as unit
			FROM 
				prescriptions_temp pt
				JOIN prescriptions_temp_detail ptd ON (ptd.prescription_id = pt.id)
				JOIN ref_drugs rd ON (rd.code = ptd.drug_code)
				JOIN ref_paramedics rp ON (rp.id = pt.doctor_id)
			WHERE 
				pt.visit_id=:visit_id 
				AND rd.satusehat_kfa_code IS NOT NULL 
				AND rd.satusehat_kfa_name IS NOT NULL
			ORDER BY ptd.id
    `

	GetMedicationDispenseByVisitId = `
			SELECT
				pt.visit_id as visit_id,
				pt.jenis_pasien as patient_type,
				pt.date as date,
				rd.code as drug_code,
				pt.id as prescription_id,
				ptd.id as prescription_detail_id,
				rd.satusehat_kfa_code as kfa_code,
				rd.satusehat_kfa_name as kfa_name,
				ptd.jenis as type,
				rp.satusehat_practitioner_id as practitioner_id,
				rp.name as paramedic_name,
				dtd.batch_number as batch_number,
				dtd.expired_date as expired_date,
				v.prescription_start_date as prescription_start_date,
				v.drug_received_by_patient_date as drug_received_date
			FROM 
				prescriptions pt
				JOIN prescriptions_detail ptd ON (ptd.prescription_id = pt.id)
				JOIN ref_drugs rd ON (rd.code = ptd.drug_code)
				JOIN ref_paramedics rp ON (rp.id = pt.doctor_id)
				JOIN drugs_transaction_detail dtd ON (dtd.prescription_detail_id = ptd.id)
				JOIN visits v ON (v.id = pt.visit_id)
			WHERE 
				pt.visit_id=:visit_id 
			ORDER BY ptd.id
    `
)

type SahabatQuery struct {
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
	queryOps := &SahabatQuery{
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

func (f *SahabatQuery) GetVisitBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]entity.Visit, error) {
	parameter := map[string]any{
		"start_date": startDate,
		"end_date":   endDate,
	}

	var results []entity.Visit

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

func (f *SahabatQuery) GetDiagnosisByVisitId(ctx context.Context, visitId string) ([]entity.Diagnosis, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []entity.Diagnosis

	err := f.getDiagnosisByVisitStmt.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *SahabatQuery) GetMedicationRequestByVisitId(ctx context.Context, visitId string) ([]entity.MedicationRequest, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []entity.MedicationRequest

	err := f.getMedicationRequestByVisitStmt.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *SahabatQuery) GetMedicationDispenseByVisitId(ctx context.Context, visitId string) ([]entity.MedicationDispense, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []entity.MedicationDispense

	err := f.getMedicationDispenseByVisitStmt.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *SahabatQuery) GetProcedureByVisitId(ctx context.Context, visitId string) ([]entity.Procedure, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []entity.Procedure

	err := f.getProcedureByVisitStmt.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *SahabatQuery) GetObservationLabByVisitId(ctx context.Context, visitId string) ([]entity.ObservationLab, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []entity.ObservationLab

	err := f.getObservationLabByVisitId.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *SahabatQuery) GetObservationRadiologyByVisitId(ctx context.Context, visitId string) ([]entity.ObservationRadiology, error) {
	parameter := map[string]any{
		"visit_id": visitId,
	}

	var results []entity.ObservationRadiology

	err := f.getObservationLabByVisitId.SelectContext(ctx, &results, parameter)

	if err != nil {
		return nil, err
	}

	return results, nil
}
