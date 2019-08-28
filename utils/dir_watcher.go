package utils

import (
	"time"
)

type DirWatcher struct {
	watchedDirs  []string
	dirChecksum  string
	ChangeChan   chan bool
	pollCooldown time.Duration
}

func MakeDirWatcher(watchedDirs []string, pollCooldown time.Duration) DirWatcher {
	watcher := DirWatcher{
		watchedDirs:  watchedDirs,
		ChangeChan:   make(chan bool),
		pollCooldown: pollCooldown,
	}
	return watcher
}

func (w *DirWatcher) Watch() {
	go func() {
		for {
			time.Sleep(w.pollCooldown)
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
