package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const USAGE = `usage: nginx-reloader [--cooldown SECONDS] [--watch DIR [DIR ...]] [--nginx-command NGINX_EXECUTABLE [NGINX_EXECUTABLE_OPTION [NGINX_EXECUTABLE_OPTION ...]]]
options:
--cooldown	
	seconds to wait after each reload
	default: 3
--watch
	space-separated directories to watch
	default: /etc/nginx/conf.d
--nginx-command
	command to start nginx with
	default: nginx -g "daemon off;"
`

const DEFAULT_POLL_COOLDOWN = 3 * time.Second

var DEFAULT_DIRS = []string{"/etc/nginx/conf.d"}
var DEFAULT_NGINX_COMMAND = []string{"nginx", "-g", "daemon off;"}

func ParseOptions(args []string) (pollCooldown time.Duration, watchedDirs []string, nginxCommand []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	parser := argParser{}
	pollCooldown, watchedDirs, nginxCommand = parser.parse(args)
	return pollCooldown, watchedDirs, nginxCommand, err
}

type argParser struct {
	cooldown     time.Duration
	dirs         []string
	nginxCommand []string

	parsedCooldown     bool
	parsedWatch        bool
	parsedNginxCommand bool

	args []string
}

func (p *argParser) parse(args []string) (cooldown time.Duration, dirs []string, nginxCommand []string) {
	switch len(args) {
	case 1:
	case 2:
		Panicf(USAGE)
	default:
		p.args = args[1:]
		p.parseStart()
	}

	if !p.parsedCooldown {
		p.cooldown = DEFAULT_POLL_COOLDOWN
	}
	if !p.parsedWatch {
		p.dirs = DEFAULT_DIRS
	}
	if !p.parsedNginxCommand {
		p.nginxCommand = DEFAULT_NGINX_COMMAND
	}

	return p.cooldown, p.dirs, p.nginxCommand
}

func (p *argParser) parseStart() {
	if len(p.args) == 0 {
		return
	}
	switch p.args[0] {
	case "--watch":
		p.parseWatch()
	case "--cooldown":
		p.parseCooldown()
	case "--nginx-command":
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
		case "--cooldown":
			p.args = p.args[idx+2:]
			p.parseCooldown()
			return
		case "--watch":
			p.args = p.args[idx+2:]
			p.parseWatch()
			return
		case "--nginx-command":
			p.args = p.args[idx+2:]
			p.parseNginxCommand()
			return
		default:
			p.dirs = append(p.dirs, el)
		}
	}
}

func (p *argParser) parseCooldown() {
	if p.parsedCooldown {
		Panicf("duplicate '--cooldown' option\n" + USAGE)
	}
	p.parsedCooldown = true
	if len(p.args) < 2 {
		Panicf("empty '--cooldown' option\n" + USAGE)
	}
	cooldown, err := strconv.Atoi(p.args[1])
	if err != nil {
		Panicf("invalid value for '--cooldown' option, expected an integer, got '%v'", p.args[1])
	}
	if cooldown < 0 {
		Panicf("watch cooldown must be >= 0, got '%v'", cooldown)
	}
	p.cooldown = time.Duration(cooldown) * time.Second
	p.args = p.args[2:]
	p.parseStart()
}

func (p *argParser) parseNginxCommand() {
	if len(p.args) < 2 {
		Panicf("empty command after '--nginx-command' option\n" + USAGE)
	}
	p.parsedNginxCommand = true
	for _, el := range p.args[1:] {
		p.nginxCommand = append(p.nginxCommand, el)
	}
}
