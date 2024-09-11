//go:build sleman

package simrs

import (
	"context"
	"github.com/jasoet/fhir-worker/pkg/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueryOps_Fetch(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	config := &db.ConnectionConfig{
		DbType:       db.MSSQL,
		Host:         "202.162.34.148",
		Port:         1433,
		Username:     "sa",
		Password:     "indomieseleraku",
		DbName:       "rsud_test",
		Timeout:      3 * time.Second,
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	}

	pool, err := config.Pool()
	assert.NoError(t, err)

	queryOps, err := NewQuery(pool)
	assert.NoError(t, err)
	assert.NotNil(t, queryOps)

	startDate := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, time.June, 31, 0, 0, 0, 0, time.UTC)

	visits, err := queryOps.GetVisitBetween(ctx, startDate, endDate)
	assert.NoError(t, err)
	assert.NotNil(t, visits)
}
