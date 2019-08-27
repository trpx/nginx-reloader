package utils

import (
	"time"
)

type DirWatcher struct {
	WatchedDirs  []string
	dirChecksum  string
	ChangeChan   chan bool
	PollInterval time.Duration
}

func MakeDirWatcher(WatchedDirs []string, PollInterval time.Duration) DirWatcher {
	watcher := DirWatcher{
		WatchedDirs:  WatchedDirs,
		ChangeChan:   make(chan bool),
		PollInterval: PollInterval,
	}
	return watcher
}

func (w *DirWatcher) Watch() {
	go func() {
		for {
			time.Sleep(w.PollInterval)
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
	for _, watchedDir := range w.WatchedDirs {
		w.dirChecksum += calcDirChecksum(watchedDir) + ";"
	}
}

func (w *DirWatcher) onChange() {
	w.ChangeChan <- true
}
