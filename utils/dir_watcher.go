package utils

import (
	"time"
)

type DirWatcher struct {
	watchedDirs  []string
	dirChecksum  string
	ChangeChan   chan bool
	pollInterval time.Duration
}

func MakeDirWatcher(watchedDirs []string, pollInterval time.Duration) DirWatcher {
	watcher := DirWatcher{
		watchedDirs:  watchedDirs,
		ChangeChan:   make(chan bool),
		pollInterval: pollInterval,
	}
	return watcher
}

func (w *DirWatcher) Watch() {
	go func() {
		for {
			time.Sleep(w.pollInterval)
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
		w.dirChecksum += calcDirChecksum(watchedDir) + ";"
	}
}

func (w *DirWatcher) onChange() {
	w.ChangeChan <- true
}
