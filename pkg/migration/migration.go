package migration

import (
	"ariga.io/atlas-go-sdk/atlasexec"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	Path          string
	NeedMigration bool
}

func Do(logger *zap.Logger, migrationConfig Config, dbConnStr string) error {
	if !migrationConfig.NeedMigration {
		return nil
	}

	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS(migrationConfig.Path),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to load working directory")
	}

	defer workdir.Close()

	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		return errors.Wrap(err, "failed to initialize client")
	}

	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: dbConnStr,
	})
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}

	logger.Info("success migrations applied", zap.Any("applied files", res.Applied))

	return nil
}
