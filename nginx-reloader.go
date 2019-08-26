package main

import (
	"github.com/trpx/nginx-reloader/utils"
	"os"
	"time"
)

func main() {

	pollEvery, watchedDirs, nginxOptions := utils.ParseDirsOptions()

	watcher := utils.DirWatcher{
		WatchedDirs: watchedDirs,
		PollEvery:   time.Duration(pollEvery) * time.Second,
		ChangeChan:  make(chan bool),
	}

	watcher.CalcChecksum()

	NginxRunner := utils.NginxRunner{
		SignalsChan:  make(chan os.Signal),
		ChangeChan:   watcher.ChangeChan,
		NginxOptions: nginxOptions,
	}

	cmd := NginxRunner.StartNginx()

	watcher.Watch()

	err := cmd.Wait()
	if err != nil {
		utils.Fatalf("nginx process encountered an error during execution:\n%v\n", err)
	}
}
