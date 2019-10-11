package main

import (
	"errors"
	"flag"
	"testing"
)

const (
	initErrorStr = "gherror must be initialized first with Register()"
	testSlug     = "hazcod/gherror"
)

// go test -token=<github-token-here>
var githubToken = flag.String("token", "", "github auth token")

func TestRegister(t *testing.T) {
	if err := Report(errors.New("foo"), nil); err == nil || err.Error() != initErrorStr {
		t.Errorf("report: %+v", err)
	}

	if err := Register("xxx", testSlug); err != nil {
		t.Errorf("register: %+v", err)
	}

	if err := Report(errors.New("foo"), nil); err != nil && err == errors.New(initErrorStr) {
		t.Errorf("report: %+v", err)
	}
}

func TestReport(t *testing.T) {
	if err := Register(*githubToken, testSlug); err != nil {
		t.Errorf("register: %+v", err)
	}

	if err := Report(errors.New("test error"), nil); err != nil {
		t.Errorf("report: %+v", err)
	}
}
