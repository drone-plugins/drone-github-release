// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	"fmt"
	"os"

	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/urfave/cli/v2"

	"github.com/drone-plugins/drone-github-release/plugin"
)

var (
	version  = "unknown"
	settings = plugin.Settings{}
)

func main() {
	app := cli.NewApp()
	app.Name = "drone-github-release"
	app.Usage = "creates a github release"
	app.Version = version
	app.Action = run
	app.Flags = append(settingsFlags(&settings), urfave.Flags()...)

	if err := app.Run(os.Args); err != nil {
		errors.HandleExit(err)
	}
}

func run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)

	plugin := plugin.New(
		settings,
		urfave.PipelineFromContext(ctx),
		urfave.NetworkFromContext(ctx),
	)

	if err := plugin.Validate(); err != nil {
		if e, ok := err.(errors.ExitCoder); ok {
			return e
		}

		return errors.ExitMessage(fmt.Errorf("validation failed: %w", err))
	}

	if err := plugin.Execute(); err != nil {
		if e, ok := err.(errors.ExitCoder); ok {
			return e
		}

		return errors.ExitMessage(fmt.Errorf("execution failed: %w", err))
	}

	return nil
}
