package simrs

import (
	"context"
	"github.com/jasoet/fhir-worker/shared/model"
	"time"
)

type Query interface {
	GetVisitBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.Visit, error)
	GetDiagnosisByVisitId(ctx context.Context, visitId string) (model.DiagnosisList, error)
	GetMedicationRequestByVisitId(ctx context.Context, visitId string) (model.MedicationRequestList, error)
	GetMedicationDispenseByVisitId(ctx context.Context, visitId string) (model.MedicationDispenseList, error)
	GetProcedureByVisitId(ctx context.Context, visitId string) (model.ProcedureList, error)
	GetObservationLabByVisitId(ctx context.Context, visitId string) (model.ObservationLabList, error)
	GetObservationRadiologyByVisitId(ctx context.Context, visitId string) (model.ObservationRadiologyList, error)
}
