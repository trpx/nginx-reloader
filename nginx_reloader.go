package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)


// todo: add to-stdout and to-stderr logging


const NGINX_DIR = "/etc/nginx/conf.d"


func main() {

	watcher := DirWatcher{
		watchedDir: NGINX_DIR,
		pollEvery: 1 * time.Second,
	}

	watcher.calcHash()

	NginxRunner := NginxRunner{
		changeChan: watcher.changeChan,
	}

	NginxRunner.startNginx()

	watcher.watchForever()

}


type NginxRunner struct {
	nginxProc *os.Process
	signals chan os.Signal
	changeChan chan bool
}


func (r NginxRunner) startNginx() {
	r.listenSignals()
	cmd := exec.Command("nginx", "-g", "daemon", "off")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	r.nginxProc = cmd.Process
	r.forwardSignals()
}


func (r NginxRunner) listenSignals() {
	signal.Notify(
		r.signals,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT,
	)
}


func (r NginxRunner) forwardSignals() {
	// Forward signals to nginx process
	go func() {
		for sig := range r.signals {
			err := r.nginxProc.Signal(sig)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}


func (r NginxRunner) reloadOnChange() {
	go func() {
		for {
			<- r.changeChan
			r.reloadNginx()
		}
	}()
}


func (r NginxRunner) reloadNginx() {
	err := r.nginxProc.Signal(syscall.SIGHUP)
	if err != nil {
		log.Fatal(err)
	}
}


type DirWatcher struct {
	watchedDir string
	dirHash string
	changeChan chan bool
	pollEvery time.Duration
}


func (w DirWatcher) watchForever() {
	for {
		time.Sleep(w.pollEvery)
		knownHash := w.dirHash
		w.calcHash()
		hasChanged := knownHash != w.dirHash
		if hasChanged {
			w.onChange()
		}
	}
}


func (w DirWatcher) calcHash() {
	w.dirHash = calcDirHash(w.watchedDir)
}


func (w DirWatcher) onChange() {
	w.changeChan <- true
}


func calcDirHash(dir string) (dirHash string) {

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fileInfos {
		if ! fi.IsDir() {
			dirHash += fmt.Sprintf("%s %d\n", fi.Name(), fi.ModTime().Unix())
		}
	}

	return dirHash
}
