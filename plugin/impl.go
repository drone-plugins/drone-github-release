// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

// Settings for the plugin.
type Settings struct {
	GitHubURL            string
	APIKey               string
	Files                cli.StringSlice
	FileExists           string
	Checksum             cli.StringSlice
	ChecksumFile         string
	ChecksumFlatten      bool
	Draft                bool
	Prerelease           bool
	BaseURL              string
	UploadURL            string
	Title                string
	Note                 string
	Overwrite            bool
	GenerateReleaseNotes bool

	baseURL   *url.URL
	uploadURL *url.URL
	uploads   []string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	var err error

	if p.pipeline.Build.Event != "tag" {
		return fmt.Errorf("github release plugin is only available for tags")
	}

	if p.settings.APIKey == "" {
		return fmt.Errorf("no api key provided")
	}

	if !fileExistsValues[p.settings.FileExists] {
		return fmt.Errorf("invalid value for file_exists")
	}

	if p.settings.BaseURL != "" && p.settings.UploadURL != "" {
		fmt.Printf("Both base_url and upload_url are deprecated. Please remove them from your config!")

		if !strings.HasSuffix(p.settings.BaseURL, "/") {
			p.settings.BaseURL = p.settings.BaseURL + "/"
		}
		p.settings.baseURL, err = url.Parse(p.settings.BaseURL)
		if err != nil {
			return fmt.Errorf("failed to parse base url: %w", err)
		}

		if !strings.HasSuffix(p.settings.UploadURL, "/") {
			p.settings.UploadURL = p.settings.UploadURL + "/"
		}
		p.settings.uploadURL, err = url.Parse(p.settings.UploadURL)
		if err != nil {
			return fmt.Errorf("failed to parse upload url: %w", err)
		}
	} else {
		p.settings.baseURL, p.settings.uploadURL, err = gitHubURLs(p.settings.GitHubURL)
		if err != nil {
			return fmt.Errorf("failed to get GitHub urls: %w", err)
		}
	}

	if p.settings.Note != "" {
		if p.settings.Note, err = readStringOrFile(p.settings.Note); err != nil {
			return fmt.Errorf("error while reading %s: %w", p.settings.Note, err)
		}
	}

	if p.settings.Title != "" {
		if p.settings.Title, err = readStringOrFile(p.settings.Title); err != nil {
			return fmt.Errorf("error while reading %s: %w", p.settings.Note, err)
		}
	}

	files := p.settings.Files.Value()
	for _, glob := range files {
		globed, err := filepath.Glob(glob)

		if err != nil {
			return fmt.Errorf("failed to glob %s: %w", glob, err)
		}

		if globed != nil {
			p.settings.uploads = append(p.settings.uploads, globed...)
		}
	}

	if len(files) > 0 && len(p.settings.uploads) < 1 {
		return fmt.Errorf("failed to find any file to release")
	}

	checksum := p.settings.Checksum.Value()
	if len(checksum) > 0 {
		p.settings.uploads, err = writeChecksums(p.settings.uploads, checksum, p.settings.ChecksumFile, p.settings.ChecksumFlatten)

		if err != nil {
			return fmt.Errorf("failed to write checksums: %w", err)
		}
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.settings.APIKey})
	tc := oauth2.NewClient(
		context.WithValue(context.Background(), oauth2.HTTPClient, p.network.Client),
		ts,
	)

	client := github.NewClient(tc)

	client.BaseURL = p.settings.baseURL
	client.UploadURL = p.settings.uploadURL

	rc := releaseClient{
		Client:               client,
		Context:              p.network.Context,
		Owner:                p.pipeline.Repo.Owner,
		Repo:                 p.pipeline.Repo.Name,
		Tag:                  strings.TrimPrefix(p.pipeline.Commit.Ref, "refs/tags/"),
		Draft:                p.settings.Draft,
		Prerelease:           p.settings.Prerelease,
		FileExists:           p.settings.FileExists,
		Title:                p.settings.Title,
		Note:                 p.settings.Note,
		Overwrite:            p.settings.Overwrite,
		GenerateReleaseNotes: p.settings.GenerateReleaseNotes,
	}

	release, err := rc.buildRelease()

	if err != nil {
		return fmt.Errorf("failed to create the release: %w", err)
	}

	if err := rc.uploadFiles(*release.ID, p.settings.uploads); err != nil {
		return fmt.Errorf("failed to upload the files: %w", err)
	}

	return nil
}

func gitHubURLs(gh string) (*url.URL, *url.URL, error) {
	uri, err := url.Parse(gh)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse GitHub link")
	}

	// Remove the path in the case that DRONE_REPO_LINK was passed in
	uri.Path = ""

	if uri.Hostname() != "github.com" {
		relBaseURL, _ := url.Parse("./api/v3/")
		relUploadURL, _ := url.Parse("./api/v3/upload/")

		return uri.ResolveReference(relBaseURL), uri.ResolveReference(relUploadURL), nil
	}

	baseURL, _ := url.Parse("https://api.github.com/")
	uploadURL, _ := url.Parse("https://uploads.github.com/")

	return baseURL, uploadURL, nil
}
