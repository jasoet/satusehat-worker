package app

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/jasoet/fhir-worker/job"
	"github.com/jasoet/fhir-worker/pkg/server"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func startFunc(config *Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_log := log.With().Ctx(ctx).Str("function", "Application").Logger()

	_log.Info().
		Msg("initializing Application")

	var err error
	queryOps, err := config.Database.QueryOps()
	if err != nil {
		_log.Error().Err(err).Msg("failed to create QueryOps")
		return err
	}

	repository, err := config.Database.Repository()
	if err != nil {
		_log.Error().Err(err).Msg("failed to create Repository")
		return err
	}

	var mappingJob *job.Mapping

	mappingOptions := []job.MappingOption{
		job.WithQueryAndRepository(queryOps, repository),
	}

	if config.Mapping != nil {
		mappingOptions = append(mappingOptions, job.WithDisableConfigs(config.Mapping.DisableDiagnosis, config.Mapping.DisableLab, config.Mapping.DisableRadiology, config.Mapping.DisableProcedure, config.Mapping.DisableMedication))
		mappingOptions = append(mappingOptions, job.WithConfigDays(config.Mapping.MarkCompleteDays, config.Mapping.LastVisitDays))
	}

	mappingJob, err = job.NewMapping(mappingOptions...)
	if err != nil {
		_log.Error().Err(err).Msg("failed to create Mapping Job")
		return err
	}

	satuSehatClient := config.Satusehat.Client()
	publishOptions := []job.PublishOption{
		job.WithOrganizationId(config.Satusehat.OrganizationID),
		job.WithClientAndRepository(satuSehatClient, repository),
		job.WithConvertUtc(config.Satusehat.ConvertToUtc),
	}

	if config.Publish != nil {
		publishOptions = append(publishOptions,
			job.WithSendDelay(config.Publish.PublishDelay),
			job.WithSimulationDir(config.Publish.SimulationDir),
			job.WithSimulationMode(config.Publish.SimulationMode),
		)
	}

	publishJob, err := job.NewPublish(publishOptions...)
	if err != nil {
		_log.Error().Err(err).Msg("failed to create Publish Job")
		return err
	}

	scheduler, err := gocron.NewScheduler(gocron.WithLimitConcurrentJobs(2, gocron.LimitModeReschedule))
	if err != nil {
		_log.Error().Err(err).Msg("failed to create Gocron Scheduler ")
		return err
	}

	if !config.Job.VisitFetchDisabled {
		fetchVisitTask, err := scheduler.NewJob(
			gocron.DurationJob(
				config.Job.VisitFetchInterval,
			),
			gocron.NewTask(
				mappingJob.FetchVisit, ctx,
			),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)

		_log.Info().
			Err(err).
			Str("Fetch Visit Task Id", fetchVisitTask.ID().String()).
			Msg("initializing Mapping fetch visit task")
	} else {
		_log.Info().
			Msg("Fetch Visit task disabled")

	}

	if !config.Job.VisitFillDisabled {
		fillVisitTask, err := scheduler.NewJob(
			gocron.DurationJob(
				config.Job.VisitFillInterval,
			),
			gocron.NewTask(
				mappingJob.FillVisit, ctx,
			),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)

		_log.Info().
			Err(err).
			Str("Fill Visit Task Id", fillVisitTask.ID().String()).
			Msg("initializing mapping fill visit task")
	} else {
		_log.Info().
			Msg("Fill Visit task disabled")
	}

	if !config.Job.MarkCompleteDisabled {
		markCompleteTask, err := scheduler.NewJob(
			gocron.DurationJob(
				config.Job.MarkCompleteInterval,
			),
			gocron.NewTask(
				mappingJob.CheckComplete, ctx,
			),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)

		_log.Info().
			Err(err).
			Str("Mark Complete Task Id", markCompleteTask.ID().String()).
			Msg("initializing mapping Mark Complete task")
	} else {
		_log.Info().
			Msg("Mark Completed task disabled")
	}

	if !config.Job.PublishDisabled {
		publishTask, err := scheduler.NewJob(
			gocron.DurationJob(
				config.Job.PublishInterval,
			),
			gocron.NewTask(
				publishJob.Process, ctx,
			),
		)

		_log.Info().
			Err(err).
			Str("Publish Task Id", publishTask.ID().String()).
			Msg("initializing publish task")
	} else {
		_log.Info().
			Msg("Publish task disabled")
	}

	server.Start(config.Port,
		func(e *echo.Echo) {
			log.Info().
				Msgf("starting scheduler")

			scheduler.Start()
		},
		func(e *echo.Echo) {
			log.Info().
				Msgf("canceling jobs, stopping scheduler, closing database connection..")

			cancel()
			_ = scheduler.Shutdown()

			log.Info().
				Msgf("scheduler stopped")
		},
	)

	return nil
}

var cfgFile string
var debugMode bool

//go:embed config_example.yaml
var configExample string

func NewCommand() *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show the example configuration",
		Long:  `This will show example configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(configExample)
		},
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the application",
		Long:  `This command starts the worker`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().
				Str("config", cfgFile).
				Msgf("loading config file")

			config, err := loadConfig(cfgFile)
			if err != nil {
				log.Error().Err(err).Msg("config file invalid")
				return err
			}

			ctx := context.WithValue(context.Background(), "config", config)
			cmd.SetContext(ctx)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, ok := cmd.Context().Value("config").(*Config)
			if !ok || config == nil {
				log.Error().
					Str("config_type", fmt.Sprintf("%T", config)).
					Msg("config file invalid")
				return fmt.Errorf("config file invalid")
			}

			return startFunc(config)
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "fhir-worker",
		Short: "CLI application to sync data to SatuSehat",
		Long:  `CLI application that periodically fetch data from database and push to SatuSehat`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			if debugMode {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "<empty config>", "config file path, use `config` command to see the example")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug mode")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(startCmd)

	return rootCmd
}
