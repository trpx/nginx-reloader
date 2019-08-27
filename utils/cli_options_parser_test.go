package utils

import (
	"testing"
	"time"
)

type TestCase struct {
	args         []string
	pollInterval int
	watchedDirs  []string
	nginxOptions []string
	err          error
}

var testCaseSuit = []TestCase{
	{
		[]string{"nginx-reloader", "--watch", "/"},
		3,
		[]string{"/"},
		[]string{"-g", "daemon off;"},
		nil,
	},
}

func TestParseOptions(t *testing.T) {
	for _, expected := range testCaseSuit {

		pollInterval, watchedDirs, nginxOptions, err := ParseOptions(expected.args)
		expectedPollInterval := time.Duration(expected.pollInterval) * time.Second
		if pollInterval != expectedPollInterval {
			t.Errorf("Expected %#v pollInterval, got %#v", expected.pollInterval, pollInterval)
		}

		watchedDirsErr := false
		if len(watchedDirs) != len(expected.watchedDirs) {
			watchedDirsErr = true
		} else {
			for idx, dir := range watchedDirs {
				if dir != expected.watchedDirs[idx] {
					watchedDirsErr = true
				}
			}
		}
		if watchedDirsErr {
			t.Errorf("Expected %#v watchedDirs, got %#v", expected.watchedDirs, watchedDirs)
		}

		nginxOptionsErr := false
		if len(nginxOptions) != len(expected.nginxOptions) {
			nginxOptionsErr = true
		} else {
			for idx, dir := range nginxOptions {
				if dir != expected.nginxOptions[idx] {
					nginxOptionsErr = true
				}
			}
		}
		if nginxOptionsErr {
			t.Errorf("Expected %#v nginxOptions, got %#v", expected.nginxOptions, nginxOptions)
		}

		if err != expected.err {
			t.Errorf("Expected %#v err, got %#v", expected.err, err)
		}
	}
}
