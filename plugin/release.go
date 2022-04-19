// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/google/go-github/v41/github"
)

// Release holds ties the drone env data and github client together.
type releaseClient struct {
	*github.Client
	context.Context
	Owner                string
	Repo                 string
	Tag                  string
	Draft                bool
	Prerelease           bool
	FileExists           string
	Title                string
	Note                 string
	Overwrite            bool
	GenerateReleaseNotes bool
}

func (rc *releaseClient) buildRelease() (*github.RepositoryRelease, error) {
	// first attempt to get a release by that tag
	release, err := rc.getRelease()

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a release: %w", err)
	}

	if release == nil {
		// if no release was found by that tag, create a new one
		release, err = rc.newRelease()
	} else {
		// update release if exists
		release, err = rc.editRelease(*release)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create or edit a release: %w", err)
	}

	return release, nil
}

func (rc *releaseClient) getRelease() (*github.RepositoryRelease, error) {

	listOpts := &github.ListOptions{PerPage: 10}

	for {
		// get list of releases (10 releases per page)
		releases, resp, err := rc.Client.Repositories.ListReleases(rc.Context, rc.Owner, rc.Repo, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}

		// browse through current release page
		for _, release := range releases {

			// return release associated to the given tag (can only be one)
			if release.GetTagName() == rc.Tag {
				fmt.Printf("Found release %d for tag %s\n", release.GetID(), release.GetTagName())
				return release, nil
			}
		}

		// end of list found without finding a matching release
		if resp.NextPage == 0 {
			fmt.Println("no existing release (draft) found for the given tag")
			return nil, nil
		}

		// go to next page in the next iteration
		listOpts.Page = resp.NextPage
	}
}

func (rc *releaseClient) editRelease(targetRelease github.RepositoryRelease) (*github.RepositoryRelease, error) {
	sourceRelease := &github.RepositoryRelease{}

	if rc.Overwrite {
		sourceRelease.Name = &rc.Title
		sourceRelease.Body = &rc.Note
	}

	// only potentially change the draft value, if it's a draft right now
	// i.e. a drafted release will be published, but a release won't be unpublished
	if targetRelease.GetDraft() {
		fmt.Printf("DRAFT: %+v\n", rc.Draft)
		if !rc.Draft {
			fmt.Println("Publishing a release draft")
		}
		sourceRelease.Draft = &rc.Draft
	}

	modifiedRelease, _, err := rc.Client.Repositories.EditRelease(rc.Context, rc.Owner, rc.Repo, targetRelease.GetID(), sourceRelease)

	if err != nil {
		return nil, fmt.Errorf("failed to update release: %w", err)
	}

	fmt.Printf("Successfully updated %s release\n", rc.Tag)
	return modifiedRelease, nil
}

func (rc *releaseClient) newRelease() (*github.RepositoryRelease, error) {
	rr := &github.RepositoryRelease{
		TagName:              github.String(rc.Tag),
		Draft:                &rc.Draft,
		Prerelease:           &rc.Prerelease,
		Name:                 &rc.Title,
		Body:                 &rc.Note,
		GenerateReleaseNotes: &rc.GenerateReleaseNotes,
	}

	if *rr.Prerelease {
		fmt.Printf("Release %s identified as a pre-release\n", rc.Tag)
	} else {
		fmt.Printf("Release %s identified as a full release\n", rc.Tag)
	}

	if *rr.Draft {
		fmt.Printf("Release %s will be created as draft (unpublished) release\n", rc.Tag)
	} else {
		fmt.Printf("Release %s will be created and published\n", rc.Tag)
	}

	if *rr.GenerateReleaseNotes {
		fmt.Printf("Release notes for %s will be automatically generated\n", rc.Tag)
	}

	release, _, err := rc.Client.Repositories.CreateRelease(rc.Context, rc.Owner, rc.Repo, rr)

	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	fmt.Printf("Successfully created %s release\n", rc.Tag)
	return release, nil
}

func (rc *releaseClient) uploadFiles(id int64, files []string) error {
	var assets []*github.ReleaseAsset
	listOpts := &github.ListOptions{PerPage: 10}
	for {
		a, resp, err := rc.Client.Repositories.ListReleaseAssets(rc.Context, rc.Owner, rc.Repo, id, listOpts)
		if err != nil {
			return fmt.Errorf("failed to fetch existing assets: %w", err)
		}
		assets = append(assets, a...)

		// stop iteration if there is no next page
		if resp.NextPage == 0 {
			break
		}

		listOpts.Page = resp.NextPage
	}

	var uploadFiles []string

files:
	for _, file := range files {
		for _, asset := range assets {
			if *asset.Name == path.Base(file) {
				switch rc.FileExists {
				case "overwrite":
					// do nothing
				case "fail":
					return fmt.Errorf("asset file %s already exists", path.Base(file))
				case "skip":
					fmt.Printf("Skipping pre-existing %s artifact\n", *asset.Name)
					continue files
				default:
					return fmt.Errorf("internal error, unknown file_exist value %s", rc.FileExists)
				}
			}
		}

		uploadFiles = append(uploadFiles, file)
	}

	for _, file := range uploadFiles {
		handle, err := os.Open(file)

		if err != nil {
			return fmt.Errorf("failed to read %s artifact: %w", file, err)
		}

		for _, asset := range assets {
			if *asset.Name == path.Base(file) {
				if _, err := rc.Client.Repositories.DeleteReleaseAsset(rc.Context, rc.Owner, rc.Repo, *asset.ID); err != nil {
					return fmt.Errorf("failed to delete %s artifact: %w", file, err)
				}

				fmt.Printf("Successfully deleted old %s artifact\n", *asset.Name)
			}
		}

		uo := &github.UploadOptions{Name: path.Base(file)}

		if _, _, err = rc.Client.Repositories.UploadReleaseAsset(rc.Context, rc.Owner, rc.Repo, id, uo, handle); err != nil {
			return fmt.Errorf("failed to upload %s artifact: %w", file, err)
		}

		fmt.Printf("Successfully uploaded %s artifact\n", file)
	}

	return nil
}
