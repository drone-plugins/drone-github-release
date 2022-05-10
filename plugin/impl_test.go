// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"testing"
)

func TestValidate(t *testing.T) {
	t.Skip()
}

func TestExecute(t *testing.T) {
	t.Skip()
}

func TestGitHubURLs(t *testing.T) {
	// GitHub case
	actualBaseURL, actualUploadURL, _ := gitHubURLs("https://github.com/drone-plugins/drone-release-download")
	expectedBaseURL := "https://api.github.com/"
	if actualBaseURL.String() != expectedBaseURL {
		t.Errorf("Unexpected base API URL (Got: %s, Expected: %s", actualBaseURL.String(), expectedBaseURL)
	}
	expectedUploadURL := "https://uploads.github.com/"
	if actualUploadURL.String() != expectedUploadURL {
		t.Errorf("Unexpected upload API URL (Got: %s, Expected: %s", actualUploadURL.String(), expectedUploadURL)
	}

	// GitHub Enterprise case
	actualBaseURL, actualUploadURL, _ = gitHubURLs("https://github.enterprise.drone.io/drone-plugins/drone-release-download")
	expectedBaseURL = "https://github.enterprise.drone.io/api/v3/"
	if actualBaseURL.String() != expectedBaseURL {
		t.Errorf("Unexpected base API URL (Got: %s, Expected: %s", actualBaseURL.String(), expectedBaseURL)
	}
	expectedUploadURL = "https://github.enterprise.drone.io/api/v3/upload/"
	if actualUploadURL.String() != expectedUploadURL {
		t.Errorf("Unexpected upload API URL (Got: %s, Expected: %s", actualUploadURL.String(), expectedUploadURL)
	}
}
