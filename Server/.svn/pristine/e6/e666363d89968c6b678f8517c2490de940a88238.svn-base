package Framework

import (
	"context"
	"github.com/hsgames/gold/app"
	"log/slog"
)

type options struct {
	Init     []func() error
	Destroy  []func() error
	Services []func() (app.Service, error)
}

type Option func(opts *options)

func WithInit(f ...func() error) Option {
	return func(opts *options) {
		opts.Init = append(opts.Init, f...)
	}
}

func WithDestroy(f ...func() error) Option {
	return func(opts *options) {
		opts.Destroy = append(opts.Destroy, f...)
	}
}

func WithServices(fs ...func() (app.Service, error)) Option {
	return func(opts *options) {
		opts.Services = append(opts.Services, fs...)
	}
}

func Run(opts ...Option) {
	var (
		err      error
		service  app.Service
		services []app.Service
		appOpts  []app.Option
		opt      options
	)

	for _, v := range opts {
		v(&opt)
	}

	if opt.Init != nil {
		for _, v := range opt.Init {
			if err = v(); err != nil {
				slog.Error("framework init", slog.Any("error", err))
				return
			}
		}
		slog.Info("framework init ok")
	}

	defer func() {
		if err == nil {
			slog.Info("framework exit ok")
		}
	}()

	if opt.Destroy != nil {
		defer func() {
			for _, v := range opt.Destroy {
				if err = v(); err != nil {
					slog.Error("framework destroy", slog.Any("error", err))
					return
				}
				slog.Info("framework destroy ok")
			}
		}()
	}

	for _, v := range opt.Services {
		if service, err = v(); err != nil {
			slog.Error("framework new service", slog.Any("error", err))
		}
		services = append(services, service)
	}

	appOpts = append(appOpts, app.WithServices(services...))

	if err = app.New(appOpts...).Run(context.Background()); err != nil {
		slog.Error("framework exit", slog.Any("error", err))
		return
	}
}
