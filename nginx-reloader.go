package main

import (
	"errors"
	"fmt"
	"github.com/trpx/nginx-reloader/utils"
	"os"
	"time"
)

// CLI
func main() {

	pollCooldown, watchedDirs, nginxCommand, err := utils.ParseOptions(os.Args)
	if err != nil {
		utils.Fatalf("%v", err)
	}

	err = StartNginxReloader(pollCooldown, watchedDirs, nginxCommand)

	if err != nil {
		utils.Fatalf("%v", err)
	}
}

// Programmatic API
func StartNginxReloader(pollCooldown time.Duration, watchedDirs []string, nginxCommand []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	validateWatchedDirs(watchedDirs)

	watcher := utils.MakeDirWatcher(watchedDirs, pollCooldown)

	watcher.CalcChecksum()

	NginxRunner := utils.MakeNginxRunner(watcher.ChangeChan, nginxCommand)

	cmd := NginxRunner.StartNginx()

	watcher.Watch()

	err = cmd.Wait()
	if err != nil {
		utils.Panicf("nginx process encountered an error during execution:\n%v\n", err)
	}
	return err
}

func validateWatchedDirs(watchedDirs []string) {
	for _, el := range watchedDirs {
		stat, err := os.Stat(el)
		if err != nil {
			if os.IsNotExist(err) {
				utils.Panicf("watched directory '%v' does not exist", el)
			} else {
				utils.Panicf("couldn't stat watched directory '%v'", el)
			}
		}
		if !stat.IsDir() {
			utils.Panicf("watched path '%v' is not a directory", el)
		}
	}
}
