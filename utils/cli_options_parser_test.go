package utils

import (
	"testing"
	"time"
)

type TestCase struct {
	args         []string
	pollInterval int
	watchedDirs  []string
	nginxCommand []string
	err          bool
}

var testCaseSuit = []TestCase{
	{
		[]string{"nginx-reloader", "--watch", "/"},
		3,
		[]string{"/"},
		[]string{"nginx", "-g", "daemon off;"},
		false,
	},
	{
		[]string{"nginx-reloader"},
		3,
		[]string{"/etc/nginx/conf.d"},
		[]string{"nginx", "-g", "daemon off;"},
		false,
	},
	{
		[]string{"nginx-reloader", "--watch", "/", ".", "--interval", "1", "--", "nginx-entrypoint.sh", "-g", "compression off;"},
		1,
		[]string{"/", "."},
		[]string{"nginx-entrypoint.sh", "-g", "compression off;"},
		false,
	},
	{
		[]string{"nginx-reloader", "--interval", "1"},
		1,
		[]string{"/etc/nginx/conf.d"},
		[]string{"nginx", "-g", "daemon off;"},
		false,
	},
	// negative interval
	{
		args: []string{"nginx-reloader", "--interval", "-1"},
		err:  true,
	},
	// unknown option
	{
		args: []string{"nginx-reloader", "--interval", "1", "--unknown-option"},
		err:  true,
	},
	// empty nginx entrypoint option '--'
	{
		args: []string{"nginx-reloader", "--interval", "1", "--"},
		err:  true,
	},
	// duplicate --interval option
	{
		args: []string{"nginx-reloader", "--interval", "1", "--interval", "1"},
		err:  true,
	},
	// duplicate --watch option
	{
		args: []string{"nginx-reloader", "--watch", "/", "--watch", "/"},
		err:  true,
	},
	// empty --interval option
	{
		args: []string{"nginx-reloader", "--interval"},
		err:  true,
	},
	// unknown option
	{
		args: []string{"nginx-reloader", "--unknown-option"},
		err:  true,
	},
	// empty nginx entrypoint option
	{
		args: []string{"nginx-reloader", "--"},
		err:  true,
	},
}

func TestParseOptions(t *testing.T) {
	for _, expected := range testCaseSuit {

		pollInterval, watchedDirs, nginxCommand, err := ParseOptions(expected.args)

		if (err != nil) != expected.err {
			if err != nil {
				t.Errorf(
					"Got unexpected error %#v when parsing args %#v",
					err, expected.args,
				)
			} else {
				t.Errorf(
					"Haven't got expected error when parsing args %#v",
					expected.args,
				)
			}
			continue
		}

		if err != nil {
			continue
		}

		expectedPollInterval := time.Duration(expected.pollInterval) * time.Second
		if pollInterval != expectedPollInterval {
			t.Errorf(
				"Expected %s pollInterval, got %s when parsing args %#v",
				expectedPollInterval, pollInterval, expected.args,
			)
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
			t.Errorf(
				"Expected %#v watchedDirs, got %#v when parsing args %#v",
				expected.watchedDirs, watchedDirs, expected.args,
			)
		}

		nginxCommandErr := false
		if len(nginxCommand) != len(expected.nginxCommand) {
			nginxCommandErr = true
		} else {
			for idx, dir := range nginxCommand {
				if dir != expected.nginxCommand[idx] {
					nginxCommandErr = true
				}
			}
		}
		if nginxCommandErr {
			t.Errorf(
				"Expected %#v nginxCommand, got %#v when parsing args %#v",
				expected.nginxCommand, nginxCommand, expected.args,
			)
		}
	}
}
