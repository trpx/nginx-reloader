package utils

import (
	"time"
)

type DirWatcher struct {
	WatchedDirs []string
	dirChecksum string
	ChangeChan  chan bool
	PollEvery   time.Duration
}

func (w *DirWatcher) Watch() {
	go func() {
		for {
			time.Sleep(w.PollEvery)
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
