package app

import (
	internalDb "github.com/jasoet/fhir-worker/internal/db"
	"github.com/jasoet/fhir-worker/internal/satusehat"
	"github.com/jasoet/fhir-worker/pkg/db"
	"github.com/jasoet/fhir-worker/pkg/util"
	"github.com/jasoet/fhir-worker/simrs"
	"github.com/spf13/viper"
	"time"
)

func loadConfig(configFile string) (*Config, error) {
	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")

	viperConfig.SetConfigFile(configFile)
	err := viperConfig.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config

	err = viperConfig.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type MappingConfig struct {
	MarkCompleteDays  int  `yaml:"mark_complete_days" mapstructure:"mark_complete_days"`
	LastVisitDays     int  `yaml:"last_visit_days" mapstructure:"last_visit_days"`
	DisableDiagnosis  bool `yaml:"disable_diagnosis" mapstructure:"disable_diagnosis"`
	DisableLab        bool `yaml:"disable_lab" mapstructure:"disable_lab"`
	DisableRadiology  bool `yaml:"disable_radiology" mapstructure:"disable_radiology"`
	DisableProcedure  bool `yaml:"disable_procedure" mapstructure:"disable_procedure"`
	DisableMedication bool `yaml:"disable_medication" mapstructure:"disable_medication"`
}

type DatabaseConfig struct {
	Path  *string             `yaml:"path" mapstructure:"path"`
	Paths []string            `yaml:"paths" mapstructure:"paths"`
	Simrs db.ConnectionConfig `yaml:"simrs" mapstructure:"simrs"`
}

type PublishConfig struct {
	SimulationMode bool          `yaml:"simulation_mode" mapstructure:"simulation_mode"`
	SimulationDir  string        `yaml:"simulation_dir" mapstructure:"simulation_dir"`
	PublishDelay   time.Duration `yaml:"publish_delay" mapstructure:"publish_delay"`
}

type SatuSehatConfig struct {
	ConvertToUtc   bool                  `yaml:"convert_to_utc" mapstructure:"convert_to_utc"`
	OrganizationID string                `yaml:"organization_id" mapstructure:"organization_id"`
	SatuSehat      satusehat.Credential  `yaml:"satusehat" mapstructure:"satusehat"`
	HttpClient     *satusehat.RestConfig `yaml:"http_client" mapstructure:"http_client"`
}

type JobConfig struct {
	PublishDisabled      bool          `json:"publish_disabled" mapstructure:"publish_disabled"`
	VisitFetchDisabled   bool          `json:"visit_fetch_disabled" mapstructure:"visit_fetch_disabled"`
	VisitFillDisabled    bool          `json:"visit_fill_disabled" mapstructure:"visit_fill_disabled"`
	MarkCompleteDisabled bool          `json:"mark_complete_disabled" mapstructure:"mark_complete_disabled"`
	PublishInterval      time.Duration `yaml:"publish_interval" mapstructure:"publish_interval"`
	VisitFetchInterval   time.Duration `yaml:"visit_fetch_interval" mapstructure:"visit_fetch_interval"`
	VisitFillInterval    time.Duration `yaml:"visit_fill_interval" mapstructure:"visit_fill_interval"`
	MarkCompleteInterval time.Duration `yaml:"mark_complete_interval" mapstructure:"mark_complete_interval"`
}

type Config struct {
	Port      int             `yaml:"port" mapstructure:"port"`
	Job       JobConfig       `yaml:"job" mapstructure:"job"`
	Mapping   *MappingConfig  `yaml:"mapping" mapstructure:"mapping"`
	Publish   *PublishConfig  `yaml:"publish" mapstructure:"publish"`
	Database  DatabaseConfig  `yaml:"database" mapstructure:"database"`
	Satusehat SatuSehatConfig `yaml:"satusehat" mapstructure:"satusehat"`
}

func (s *SatuSehatConfig) Client() *satusehat.Client {
	options := []satusehat.ClientOption{
		satusehat.WithCredential(s.SatuSehat),
	}

	if s.HttpClient != nil {
		options = append(options, satusehat.WithRestConfig(*s.HttpClient))
	}

	return satusehat.NewClient(options...)
}

func (d *DatabaseConfig) QueryOps() (simrs.Query, error) {
	dbPool, err := d.Simrs.Pool()

	if err != nil {
		return nil, err
	}

	return simrs.NewQuery(dbPool)
}

func (d *DatabaseConfig) Repository() (*internalDb.Repository, error) {
	var paths []string
	if d.Path != nil && util.StringNotEmpty(*d.Path) {
		paths = append(paths, *d.Path)
	} else {
		paths = d.Paths
	}

	return internalDb.DefaultRepository(paths...)
}
