package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
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
		User     string
		Files    []string
		Password string
		BaseURL  string
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
		return fmt.Errorf("The SVN Release plugin is only available for tags")
	}

	if p.Config.User == "" {
		return fmt.Errorf("You must provide a User")
	}
	if p.Config.Password == "" {
		return fmt.Errorf("You must provide a Password")
	}
	if p.Config.BaseURL == "" {
		return fmt.Errorf("You must provide a SVN Directory URL")
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

	// if len(p.Config.Checksum) > 0 {
	// 	var (
	// 		err error
	// 	)
	//
	// 	files, err = writeChecksums(files, p.Config.Checksum)
	//
	// 	if err != nil {
	// 		return fmt.Errorf("Failed to write checksums. %s", err)
	// 	}
	// }

	// baseURL, err := url.Parse(p.Config.BaseURL)
	//
	// if err != nil {
	// 	return fmt.Errorf("Failed to parse base URL. %s", err)
	// }

	// rc := releaseClient{
	// 	User:     p.Config.User,
	// 	Password: p.Config.Password,
	// 	BaseURL:  p.Config.BaseURL,
	//
	// 	Owner: p.Repo.Owner,
	// 	Repo:  p.Repo.Name,
	// 	Tag:   filepath.Base(p.Commit.Ref),
	// 	Draft: p.Config.Draft,
	// }
	//
	// release, err := rc.buildRelease()
	// if err != nil {
	// 	return fmt.Errorf("Failed to create the release. %s", err)
	// }

	// if err := rc.uploadFiles(*release.ID, files); err != nil {
	if err := release(p.Config.User, p.Config.Password, p.Config.BaseURL, p.Config.Files); err != nil {
		return fmt.Errorf("Failed to upload the files. %s", err)
	}

	return nil
}

func release(user string, password string, url string, files []string) error {
	// bash -c 'svn co --username $$ATLAS_USER --password "$$ATLAS_PASSWORD" --depth empty --trust-server-cert --non-interactive "https://atlas.sys.comcast.net/iss/iss/x86_64/7/global" atlas'
	// @mv *.rpm atlas
	// bash -c 'cd atlas && svn add *.rpm && svn ci --trust-server-cert --non-interactive --username $$ATLAS_USER --password "$$ATLAS_PASSWORD" -m "drone-ho-b01.dna.comcast.net: sampleproject: $$BUILD_NUMBER" *.rpm'

	makeDir := exec.Command("make_svn_dir.sh", user, password, url, "svn-base-dir")
	if err := execute(makeDir); err != nil {
		return err
	}

	for _, file := range files {
		stage := exec.Command("cp", file, "svn-base-dir/")
		if err := execute(stage); err != nil {
			return fmt.Errorf("Failed to stage %s artifact: %s", file, err)
		}
	}

	push := exec.Command("push.sh", user, password, "svn-base-dir/")
	if err := execute(push); err != nil {
		return fmt.Errorf("Failed to stage artifacts: %s", err)
	}
	return nil
}
