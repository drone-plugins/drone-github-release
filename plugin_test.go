package main // Needs to be same package to be able to access the private func "readStringOrFile"

import (
	"testing"
)

func TestReadStringOrFileSelf(t *testing.T) {
	contents, err := readStringOrFile("./plugin_test.go")

	if err != nil {
		t.Error(err)
		return
	}
	if len(contents) == 0 {
		t.Errorf("Expected this file to have length > 0, was %d", len(contents))
	}
}

func TestReadStringOrFileLongString(t *testing.T) {
	s := "if the string is extremely long it will still try to ask the OS to read this as a file which in some cases will not be allowed because of the length of the file name however the plugin might try this anyways but most file systems only allow a maximum of 255 chars for a file name but up to 4096 for a full path thats a lot of characters"
	contents, err := readStringOrFile(s)

	if err != nil {
		t.Error(err)
		return
	}

	if contents != s {
		t.Error("Expected readStringOrFile to return input for a long string")
	}
}
