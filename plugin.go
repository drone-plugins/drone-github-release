package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Event string
	}

	Commit struct {
		Ref string
	}

	Config struct {
		APIKey          string
		Files           []string
		FileExists      string
		Checksum        []string
		ChecksumFile    string
		ChecksumFlatten bool
		Draft           bool
		Prerelease      bool
		BaseURL         string
		UploadURL       string
		Title           string
		Note            string
		Overwrite       bool
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
	}
)

func (p Plugin) Exec() error {
	var (
		files []string
	)

	if p.Build.Event != "tag" {
		return fmt.Errorf("The GitHub Release plugin is only available for tags")
	}

	if p.Config.APIKey == "" {
		return fmt.Errorf("You must provide an API key")
	}

	if !fileExistsValues[p.Config.FileExists] {
		return fmt.Errorf("Invalid value for file_exists")
	}

	if !strings.HasSuffix(p.Config.BaseURL, "/") {
		p.Config.BaseURL = p.Config.BaseURL + "/"
	}

	if !strings.HasSuffix(p.Config.UploadURL, "/") {
		p.Config.UploadURL = p.Config.UploadURL + "/"
	}

	var err error
	if p.Config.Note != "" {
		if p.Config.Note, err = readStringOrFile(p.Config.Note); err != nil {
			return fmt.Errorf("error while reading %s: %v", p.Config.Note, err)
		}
	}

	if p.Config.Title != "" {
		if p.Config.Title, err = readStringOrFile(p.Config.Title); err != nil {
			return fmt.Errorf("error while reading %s: %v", p.Config.Note, err)
		}
	}

	for _, glob := range p.Config.Files {
		globed, err := filepath.Glob(glob)

		if err != nil {
			return fmt.Errorf("Failed to glob %s. %s", glob, err)
		}

		if globed != nil {
			files = append(files, globed...)
		}
	}

	if len(p.Config.Checksum) > 0 {
		var (
			err error
		)

		files, err = writeChecksums(files, p.Config.Checksum, p.Config.ChecksumFile, p.Config.ChecksumFlatten)

		if err != nil {
			return fmt.Errorf("Failed to write checksums. %s", err)
		}
	}

	baseURL, err := url.Parse(p.Config.BaseURL)

	if err != nil {
		return fmt.Errorf("Failed to parse base URL. %s", err)
	}

	uploadURL, err := url.Parse(p.Config.UploadURL)

	if err != nil {
		return fmt.Errorf("Failed to parse upload URL. %s", err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Config.APIKey})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	client.BaseURL = baseURL
	client.UploadURL = uploadURL

	rc := releaseClient{
		Client:     client,
		Context:    ctx,
		Owner:      p.Repo.Owner,
		Repo:       p.Repo.Name,
		Tag:        strings.TrimPrefix(p.Commit.Ref, "refs/tags/"),
		Draft:      p.Config.Draft,
		Prerelease: p.Config.Prerelease,
		FileExists: p.Config.FileExists,
		Title:      p.Config.Title,
		Note:       p.Config.Note,
		Overwrite:  p.Config.Overwrite,
	}

	release, err := rc.buildRelease()

	if err != nil {
		return fmt.Errorf("Failed to create the release. %s", err)
	}

	if err := rc.uploadFiles(*release.ID, files); err != nil {
		return fmt.Errorf("Failed to upload the files. %s", err)
	}

	return nil
}

func readStringOrFile(input string) (string, error) {
	if len(input) > 255 {
		return input, nil
	}
	// Check if input is a file path
	if _, err := os.Stat(input); err != nil && os.IsNotExist(err) {
		// No file found => use input as result
		return input, nil
	} else if err != nil {
		return "", err
	}
	result, err := ioutil.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
