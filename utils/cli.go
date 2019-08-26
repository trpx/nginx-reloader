package utils

import (
	"os"
	"strconv"
)

const USAGE = "usage: nginx-reloader <POLL_EVERY_SECONDS> [WATCHED_DIR...] -- [NGINX_CLI_OPTION...]"

func ParseDirsOptions() (pollEvery int, dirs []string, options []string) {
	if len(os.Args) < 3 {
		// todo refactor to 0 required cli args
		//if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		//
		//}
		Panicf(USAGE)
	}

	pollEvery, err := strconv.Atoi(os.Args[1])
	if err != nil {
		Panicf(USAGE)
	}

	separatorIndex := len(os.Args)
	for idx, el := range os.Args {
		if el == "--" {
			separatorIndex = idx
			if len(os.Args) > idx+1 {
				options = os.Args[separatorIndex+1:]
			}
			break
		}
	}
	dirs = os.Args[2:separatorIndex]
	return pollEvery, dirs, options
}
