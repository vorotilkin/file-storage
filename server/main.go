package main

import (
	"context"

	"github.com/vorotilkin/file-storage/infrastructure/repositories/files"
	"github.com/vorotilkin/file-storage/interfaces"
	"github.com/vorotilkin/file-storage/pkg/configuration"
	"github.com/vorotilkin/file-storage/pkg/database"
	pkgGrpc "github.com/vorotilkin/file-storage/pkg/grpc"
	"github.com/vorotilkin/file-storage/pkg/migration"
	"github.com/vorotilkin/file-storage/pkg/s3"
	"github.com/vorotilkin/file-storage/proto"
	"github.com/vorotilkin/file-storage/usecases"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type config struct {
	Grpc struct {
		Server pkgGrpc.Config
	}
	Db        database.Config
	Migration migration.Config
	S3        s3.Config
}

func newConfig(configuration *configuration.Configuration) (*config, error) {
	c := new(config)
	err := configuration.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	opts := []fx.Option{
		fx.Provide(zap.NewProduction),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(configuration.New),
		fx.Provide(newConfig),
		fx.Provide(func(c *config) pkgGrpc.Config {
			return c.Grpc.Server
		}),
		fx.Provide(func(c *config) database.Config {
			return c.Db
		}),
		fx.Provide(database.New),
		fx.Provide(func(c *config) migration.Config { return c.Migration }),
		fx.Provide(fx.Annotate(func(c *config) string { return c.Db.PostgresDSN() }, fx.ResultTags(`name:"dsn"`))),
		fx.Provide(func(c *config) s3.Config { return c.S3 }),
		fx.Provide(fx.Annotate(pkgGrpc.NewServer,
			fx.As(new(grpc.ServiceRegistrar)),
			fx.As(new(interfaces.Hooker)))),
		fx.Provide(fx.Annotate(files.NewRepository, fx.As(new(usecases.FilesRepository)))),
		fx.Provide(fx.Annotate(s3.New, fx.As(new(usecases.S3Service)))),
		fx.Provide(fx.Annotate(usecases.NewFileStorageServer, fx.As(new(proto.FileStorageServiceServer)))),
		fx.Invoke(func(lc fx.Lifecycle, server interfaces.Hooker) {
			lc.Append(fx.Hook{
				OnStart: server.OnStart,
				OnStop:  server.OnStop,
			})
		}),
		fx.Invoke(fx.Annotate(migration.Do, fx.ParamTags("", "", `name:"dsn"`))),
		fx.Invoke(proto.RegisterFileStorageServiceServer),
	}

	app := fx.New(opts...)
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	<-app.Done()

	err = app.Stop(context.Background())
	if err != nil {
		panic(err)
	}
}
