package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "svn-release plugin"
	app.Usage = "svn-release plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "user",
			Usage:  "username to access svn url",
			EnvVar: "PLUGIN_USER,SVN_RELEASE_USER",
		},
		cli.StringSliceFlag{
			Name:   "files",
			Usage:  "list of files to upload",
			EnvVar: "PLUGIN_FILES,SVN_RELEASE_FILES",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "password to access svn url",
			EnvVar: "PLUGIN_PASSWORD,SVN_RELEASE_PASSWORD",
		},
		cli.StringFlag{
			Name:   "base-url",
			Usage:  "svn base URL to publish files to",
			EnvVar: "PLUGIN_BASE_URL,SVN_RELEASE_BASE_URL",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Event: c.String("build.event"),
		},
		Commit: Commit{
			Ref: c.String("commit.ref"),
		},
		Config: Config{
			User:     c.String("user"),
			BaseURL:  c.String("base-url"),
			Password: c.String("password"),
			Files:    c.StringSlice("files"),
		},
	}

	return plugin.Exec()
}
