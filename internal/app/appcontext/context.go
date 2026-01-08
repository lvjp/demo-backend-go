package appcontext

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"time"

	"go.lvjp.me/demo-backend-go/internal/app/config"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type AppContext struct {
	context.Context

	Input  io.Reader
	Output io.Writer
	Error  io.Writer
	Logger zerolog.Logger

	Config *config.Config
}

func NewFromCommand(cmd *cobra.Command) (*AppContext, error) {
	ret := &AppContext{
		Context: cmd.Context(),

		Input:  cmd.InOrStdin(),
		Output: cmd.OutOrStdout(),
		Error:  cmd.ErrOrStderr(),
	}

	if err := ret.init(); err != nil {
		return nil, fmt.Errorf("context initialization: %v", err)
	}

	return ret, nil
}

func (ctx *AppContext) init() error {
	if config, err := config.Load(); err != nil {
		return err
	} else {
		ctx.Config = config
	}

	if err := ctx.initLogger(); err != nil {
		return err
	}

	return nil
}

func (ctx *AppContext) initLogger() error {
	writer := ctx.Error

	var unknowFormat bool

	switch ctx.Config.Log.Format {
	case "json":
		// default is json, do nothing
	case "console":
		writer = zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
		}
	default:
		unknowFormat = true
	}

	ctx.Logger = zerolog.New(writer).With().Timestamp().Logger()
	ctx.Context = ctx.Logger.WithContext(ctx.Context)

	if unknowFormat {
		ctx.Logger.Warn().
			Str("format", ctx.Config.Log.Format).
			Msg("unknown log format, defaulting to json")
	}

	// Remove date/time flags which are already present in zerolog output
	stdlog.SetFlags(stdlog.Flags() & ^(stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds))
	stdlog.SetOutput(ctx.Logger.With().Str("level", "stdlog").Logger())

	return nil
}
