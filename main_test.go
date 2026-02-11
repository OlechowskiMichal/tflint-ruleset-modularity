package main

import "testing"

func TestVersion(t *testing.T) {
	t.Parallel()

	if version == "" {
		t.Error("version must not be empty")
	}
}
