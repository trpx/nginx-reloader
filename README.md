
# Usage

This util is meant to be used as an entrypoint in Docker containers to start
 `nginx` and reload it when `*.conf` files change.
 See the provided exemplary `Dockerfile`.

## Cli

### usage:

`nginx-reloader [--cooldown SECONDS] [--watch DIR [DIR ...]] [--nginx-command NGINX_EXECUTABLE [NGINX_EXECUTABLE_OPTION [NGINX_EXECUTABLE_OPTION ...]]]`

#### options:
`--cooldown`  	
- seconds to wait after each reload  
  default: `3`
  
`--watch`
- space-separated directories to watch  
  default: `/etc/nginx/conf.d`

`--nginx-command`
- command to start nginx with  
  default: `nginx -g "daemon off;"`


### example:

e.g. command   
`nginx-reloader --cooldown 10 --watch /etc/nginx/conf.d --nginx-command nginx -g "daemon off;"`

- starts nginx with `nginx -g "daemon off;"` command
- checksums the `*.conf` files in `/etc/nginx/conf.d` directory every `10` seconds
- reloads nginx on every change

## Other Details

### Signals
 
forwarded to nginx: SIGHUP, SIGINT, SIGTERM, SIGQUIT, SIGKILL and SIGABRT

### Logging

nginx-reloader writes to `stdout` (when changes are detected and nginx is reloaded) and `stderr` (on fatal errors which abort the execution of the program) only, which is what Docker logs expect.
