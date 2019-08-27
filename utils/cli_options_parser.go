package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const USAGE = "usage: nginx-reloader [--poll-interval SECONDS] [--watch DIR [DIR ...]] [--nginx CLI_OPTION [CLI_OPTION ...]]"
const DEFAULT_POLL_INTERVAL = 3 * time.Second

var DEFAULT_DIRS = []string{"/etc/nginx/conf.d"}
var DEFAULT_NGINX_OPTIONS []string

func ParseOptions() (pollInterval time.Duration, watchedDirs []string, nginxOptions []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	parser := argParser{}
	pollInterval, watchedDirs, nginxOptions = parser.parse()
	return pollInterval, watchedDirs, nginxOptions, err
}

type argParser struct {
	interval time.Duration
	dirs     []string
	options  []string

	parsedInterval bool
	parsedWatch    bool
	parsedOptions  bool

	args []string
}

func (p *argParser) parse() (interval time.Duration, dirs []string, options []string) {
	switch len(os.Args) {
	case 1:
		p.interval = DEFAULT_POLL_INTERVAL
		p.dirs = DEFAULT_DIRS
		p.options = DEFAULT_NGINX_OPTIONS
	case 2:
		Panicf(USAGE)
	default:
		p.args = os.Args[1:]
		p.parseStart()
	}

	return p.interval, p.dirs, p.options
}

func (p *argParser) parseStart() {
	if len(p.args) == 0 {
		return
	}
	switch p.args[0] {
	case "--watch":
		p.parseWatch()
	case "--poll-interval":
		p.parseInterval()
	case "--nginx":
		p.parseOptions()
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
		case "--poll-interval":
			p.args = p.args[idx+2:]
			p.parseInterval()
			return
		case "--nginx":
			p.args = p.args[idx+2:]
			p.parseOptions()
			return
		default:
			stat, err := os.Stat(el)
			if err != nil {
				if os.IsNotExist(err) {
					Panicf("directory '%v' does not exist", el)
				} else {
					Panicf("couldn't stat directory '%v'", el)
				}
			}
			if !stat.IsDir() {
				Panicf("'%v' is not a directory", el)
			}
			p.dirs = append(p.dirs, el)
		}
	}
}

func (p *argParser) parseInterval() {
	if p.parsedInterval {
		Panicf("duplicate '--poll-interval' option\n" + USAGE)
	}
	p.parsedInterval = true
	if len(p.args) < 2 {
		Panicf("empty '--poll-interval' option\n" + USAGE)
	}
	interval, err := strconv.Atoi(p.args[1])
	if err != nil {
		Panicf("invalid value for '--poll-interval' option, expected an integer, got '%v'", p.args[1])
	}
	if interval < 0 {
		Panicf("watch interval must be >= 0, got '%v'", interval)
	}
	p.interval = time.Duration(interval) * time.Second
	p.args = p.args[2:]
	p.parseStart()
}

func (p *argParser) parseOptions() {
	if len(p.args) < 2 {
		Panicf("empty '--nginx' option\n" + USAGE)
	}
	for _, el := range p.args[1:] {
		p.options = append(p.options, el)
	}
}
