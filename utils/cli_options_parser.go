package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const USAGE = "usage: nginx-reloader [--interval SECONDS] [--watch DIR [DIR ...]] [-- NGINX_ENTRYPOINT [NGINX OPTION [NGINX_OPTION ...]]]"

const DEFAULT_POLL_INTERVAL = 3 * time.Second

var DEFAULT_DIRS = []string{"/etc/nginx/conf.d"}
var DEFAULT_NGINX_COMMAND = []string{"nginx", "-g", "daemon off;"}

func ParseOptions(args []string) (pollInterval time.Duration, watchedDirs []string, nginxCommand []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	parser := argParser{}
	pollInterval, watchedDirs, nginxCommand = parser.parse(args)
	return pollInterval, watchedDirs, nginxCommand, err
}

type argParser struct {
	interval     time.Duration
	dirs         []string
	nginxCommand []string

	parsedInterval     bool
	parsedWatch        bool
	parsedNginxCommand bool

	args []string
}

func (p *argParser) parse(args []string) (interval time.Duration, dirs []string, nginxCommand []string) {
	switch len(args) {
	case 1:
	case 2:
		Panicf(USAGE)
	default:
		p.args = args[1:]
		p.parseStart()
	}

	if !p.parsedInterval {
		p.interval = DEFAULT_POLL_INTERVAL
	}
	if !p.parsedWatch {
		p.dirs = DEFAULT_DIRS
	}
	if !p.parsedNginxCommand {
		p.nginxCommand = DEFAULT_NGINX_COMMAND
	}

	return p.interval, p.dirs, p.nginxCommand
}

func (p *argParser) parseStart() {
	if len(p.args) == 0 {
		return
	}
	switch p.args[0] {
	case "--watch":
		p.parseWatch()
	case "--interval":
		p.parseInterval()
	case "--":
		p.parseNginxCommand()
	default:
		Panicf("unknown option '%s'\n"+USAGE, p.args[0])
	}
}

func (p *argParser) parseWatch() {
	if p.parsedWatch {
		Panicf("duplicate '--watch' option\n" + USAGE)
	}
	p.parsedWatch = true
	if len(p.args) < 2 {
		Panicf("empty '--watch' option\n" + USAGE)
	}
	p.dirs = append(p.dirs, p.args[1])
	for idx, el := range p.args[2:] {
		switch el {
		case "--interval":
			p.args = p.args[idx+2:]
			p.parseInterval()
			return
		case "--watch":
			p.args = p.args[idx+2:]
			p.parseWatch()
			return
		case "--":
			p.args = p.args[idx+2:]
			p.parseNginxCommand()
			return
		default:
			p.dirs = append(p.dirs, el)
		}
	}
}

func (p *argParser) parseInterval() {
	if p.parsedInterval {
		Panicf("duplicate '--interval' option\n" + USAGE)
	}
	p.parsedInterval = true
	if len(p.args) < 2 {
		Panicf("empty '--interval' option\n" + USAGE)
	}
	interval, err := strconv.Atoi(p.args[1])
	if err != nil {
		Panicf("invalid value for '--interval' option, expected an integer, got '%v'", p.args[1])
	}
	if interval < 0 {
		Panicf("watch interval must be >= 0, got '%v'", interval)
	}
	p.interval = time.Duration(interval) * time.Second
	p.args = p.args[2:]
	p.parseStart()
}

func (p *argParser) parseNginxCommand() {
	if len(p.args) < 2 {
		Panicf("empty command after '--' option\n" + USAGE)
	}
	p.parsedNginxCommand = true
	for _, el := range p.args[1:] {
		p.nginxCommand = append(p.nginxCommand, el)
	}
}
