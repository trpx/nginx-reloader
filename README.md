
# Alpha

Be warned: for now tested only manually, seems too be working, but automated tests are yet to come.

## Usage

This util is meant to be used as an entrypoint in Docker containers to start
 `nginx` and reload it when `*.conf` files change.
 See the provided exemplary `Dockerfile`.

## Cli

`usage: nginx-reloader [--interval SECONDS] [--watch DIR [DIR ...]] [--nginx CLI_OPTION [CLI_OPTION ...]]`

### defaults:
`nginx-reloader` with no arguments defaults to  
 `nginx-reloader --interval 3 --watch /etc/nginx/conf.d --nginx -g "daemon off;"`

### example

e.g. command   
`nginx-reloader --interval 10 --watch /etc/nginx/conf.d --nginx -g "daemon off;"`

- starts nginx
- checksums the `*.conf` files in `/etc/nginx/conf.d` directory every `10` seconds
- reloads nginx on every change

## Other Details

### Signals
 
forwarded to nginx: SIGHUP, SIGINT, SIGTERM, SIGQUIT, SIGKILL and SIGABRT

### Logging

nginx-reloader writes to `stdout` (when changes are detected and nginx is reloaded) and `stderr` (on fatal errors which abort the execution of the program) only, which is what Docker logs expect.
