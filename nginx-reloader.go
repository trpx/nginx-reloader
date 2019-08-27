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

	pollInterval, watchedDirs, nginxOptions, err := utils.ParseOptions()
	if err != nil {
		utils.Fatalf("%v", err)
	}

	fmt.Printf("%v\n%v\n%v\n", pollInterval, watchedDirs, nginxOptions)
	os.Exit(0)

	err = StartNginxReloader(pollInterval, watchedDirs, nginxOptions)

	if err != nil {
		utils.Fatalf("%v", err)
	}
}

// Programmatic API
func StartNginxReloader(pollInterval time.Duration, watchedDirs []string, nginxOptions []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	watcher := utils.MakeDirWatcher(
		watchedDirs,
		time.Duration(pollInterval)*time.Second,
	)

	watcher.CalcChecksum()

	NginxRunner := utils.MakeNginxRunner(
		watcher.ChangeChan,
		nginxOptions,
	)

	cmd := NginxRunner.StartNginx()

	watcher.Watch()

	err = cmd.Wait()
	if err != nil {
		utils.Panicf("nginx process encountered an error during execution:\n%v\n", err)
	}
	return err
}
