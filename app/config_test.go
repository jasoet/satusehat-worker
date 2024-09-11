package app

import (
	_ "embed"
	"github.com/jasoet/fhir-worker/pkg/db"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

//go:embed config_optional_example.yaml
var configWithOptionalExample string

func TestLoadConfig(t *testing.T) {
	// Write config to a temporary file
	tmpfile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(configExample)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load and validate config
	config, err := loadConfig(tmpfile.Name())
	assert.Nil(t, err, "Expected to load configuration file successfully")
	assert.Equal(t, 8000, config.Port, "Expected port to be 8000")
	assert.Equal(t, 30*time.Second, config.Job.PublishInterval, "Expected publish_interval to be 30s")
	assert.Equal(t, 10*time.Second, config.Job.VisitFetchInterval, "Expected visit_fetch_interval to be 1s")
	assert.Equal(t, 10*time.Second, config.Job.VisitFillInterval, "Expected visit_fill_interval to be 1s")
	assert.Equal(t, 10*time.Second, config.Job.MarkCompleteInterval, "Expected mark_complete_interval to be 1s")
	assert.Equal(t, db.Mysql, config.Database.Simrs.DbType, "Expected db_type to be MYSQL")
	assert.Equal(t, "internal.db", *config.Database.Path, "Expected Path to be internal.db ")
	assert.Equal(t, "localhost", config.Database.Simrs.Host, "Expected host to be localhost")
}
func TestLoadOptionalConfig(t *testing.T) {
	// Write config to a temporary file
	tmpfile, err := os.CreateTemp("", "config_optional.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(configWithOptionalExample)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load and validate config
	config, err := loadConfig(tmpfile.Name())
	assert.Nil(t, err, "Expected to load configuration file successfully")
	assert.Equal(t, 8000, config.Port, "Expected port to be 8000")
	assert.Equal(t, 30*time.Second, config.Job.PublishInterval, "Expected publish_interval to be 30s")
	assert.Equal(t, 1*time.Second, config.Job.VisitFetchInterval, "Expected visit_fetch_interval to be 1s")
	assert.Equal(t, 1*time.Second, config.Job.VisitFillInterval, "Expected visit_fill_interval to be 1s")
	assert.Equal(t, 1*time.Second, config.Job.MarkCompleteInterval, "Expected mark_complete_interval to be 1s")
	assert.Nil(t, config.Database.Path, "Expected path to be null")
	assert.Nil(t, config.Mapping, "Expected Mapping to be null")
	assert.Nil(t, config.Publish, "Expected publish to be null")
	assert.Nil(t, config.Satusehat.HttpClient, "Expected HttpClient to be null")
}
