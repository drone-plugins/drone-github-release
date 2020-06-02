// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/drone-plugins/drone-github-release/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "api-key",
			Usage:       "api key to access github api",
			EnvVars:     []string{"PLUGIN_API_KEY", "GITHUB_RELEASE_API_KEY", "GITHUB_TOKEN"},
			Destination: &settings.APIKey,
		},
		&cli.StringSliceFlag{
			Name:        "files",
			Usage:       "list of files to upload",
			EnvVars:     []string{"PLUGIN_FILES", "GITHUB_RELEASE_FILES"},
			Destination: &settings.Files,
		},
		&cli.StringFlag{
			Name:        "file-exists",
			Value:       "overwrite",
			Usage:       "what to do if file already exist",
			EnvVars:     []string{"PLUGIN_FILE_EXISTS", "GITHUB_RELEASE_FILE_EXISTS"},
			Destination: &settings.FileExists,
		},
		&cli.StringSliceFlag{
			Name:        "checksum",
			Usage:       "generate specific checksums",
			EnvVars:     []string{"PLUGIN_CHECKSUM", "GITHUB_RELEASE_CHECKSUM"},
			Destination: &settings.Checksum,
		},
		&cli.StringFlag{
			Name:        "checksum-file",
			Usage:       "name used for checksum file. \"CHECKSUM\" is replaced with the chosen method",
			EnvVars:     []string{"PLUGIN_CHECKSUM_FILE"},
			Value:       "CHECKSUMsum.txt",
			Destination: &settings.ChecksumFile,
		},
		&cli.BoolFlag{
			Name:        "checksum-flatten",
			Usage:       "include only the basename of the file in the checksum file",
			EnvVars:     []string{"PLUGIN_CHECKSUM_FLATTEN"},
			Destination: &settings.ChecksumFlatten,
		},
		&cli.BoolFlag{
			Name:        "draft",
			Usage:       "create a draft release",
			EnvVars:     []string{"PLUGIN_DRAFT", "GITHUB_RELEASE_DRAFT"},
			Destination: &settings.Draft,
		},
		&cli.BoolFlag{
			Name:        "prerelease",
			Usage:       "set the release as prerelease",
			EnvVars:     []string{"PLUGIN_PRERELEASE", "GITHUB_RELEASE_PRERELEASE"},
			Destination: &settings.Prerelease,
		},
		&cli.StringFlag{
			Name:        "base-url",
			Usage:       "api url, needs to be changed for ghe",
			Value:       "https://api.github.com/",
			EnvVars:     []string{"PLUGIN_BASE_URL", "GITHUB_RELEASE_BASE_URL"},
			Destination: &settings.BaseURL,
		},
		&cli.StringFlag{
			Name:        "upload-url",
			Usage:       "upload url, needs to be changed for ghe",
			Value:       "https://uploads.github.com/",
			EnvVars:     []string{"PLUGIN_UPLOAD_URL", "GITHUB_RELEASE_UPLOAD_URL"},
			Destination: &settings.UploadURL,
		},
		&cli.StringFlag{
			Name:        "title",
			Usage:       "file or string for the title shown in the github release",
			EnvVars:     []string{"PLUGIN_TITLE", "GITHUB_RELEASE_TITLE"},
			Destination: &settings.Title,
		},
		&cli.StringFlag{
			Name:        "note",
			Usage:       "file or string with notes for the release (example: changelog)",
			EnvVars:     []string{"PLUGIN_NOTE", "GITHUB_RELEASE_NOTE"},
			Destination: &settings.Note,
		},
		&cli.BoolFlag{
			Name:        "overwrite",
			Usage:       "force overwrite existing release informations e.g. title or note",
			EnvVars:     []string{"PLUGIN_OVERWRITE", "GITHUB_RELEASE_OVERWRIDE"},
			Destination: &settings.Overwrite,
		},
	}
}
