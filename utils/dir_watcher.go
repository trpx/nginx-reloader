package utils

import (
	"time"
)

type DirWatcher struct {
	watchedDirs     []string
	dirChecksum     string
	ChangeChan      chan bool
	dirPollInterval time.Duration
}

func MakeDirWatcher(WatchedDirs []string, PollInterval time.Duration) DirWatcher {
	watcher := DirWatcher{
		watchedDirs:     WatchedDirs,
		ChangeChan:      make(chan bool),
		dirPollInterval: PollInterval / time.Duration(len(WatchedDirs)),
	}
	return watcher
}

func (w *DirWatcher) Watch() {
	go func() {
		for {
			knownChecksum := w.dirChecksum
			w.CalcChecksum()
			hasChanged := knownChecksum != w.dirChecksum
			if hasChanged {
				w.onChange()
			}
		}
	}()
}

func (w *DirWatcher) CalcChecksum() {
	w.dirChecksum = ""
	for _, watchedDir := range w.watchedDirs {
		time.Sleep(w.dirPollInterval)
		w.dirChecksum += calcDirChecksum(watchedDir) + ";"
	}
}

func (w *DirWatcher) onChange() {
	w.ChangeChan <- true
}
