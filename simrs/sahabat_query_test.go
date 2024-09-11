//go:build !sleman

package simrs

import (
	"context"
	"github.com/jasoet/fhir-worker/pkg/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueryOps_Fetch(t *testing.T) {
	//t.Skip()
	ctx := context.Background()
	config := &db.ConnectionConfig{
		DbType:       db.Mysql,
		Host:         "103.150.191.237",
		Port:         33066,
		Username:     "simrs",
		Password:     "bismilah",
		DbName:       "simrs_old",
		Timeout:      3 * time.Second,
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	}

	pool, err := config.Pool()
	assert.NoError(t, err)

	queryOps, err := NewQuery(pool)
	assert.NoError(t, err)
	assert.NotNil(t, queryOps)

	startDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)

	visits, err := queryOps.GetVisitBetween(ctx, startDate, endDate)
	assert.NoError(t, err)
	assert.NotNil(t, visits)

	for _, visit := range visits {
		data, err := queryOps.GetDiagnosisByVisitId(ctx, visit.VisitID)
		assert.NoError(t, err)
		if data != nil {
			assert.NotEmpty(t, data)
		}
	}
}
