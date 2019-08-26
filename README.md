
# Alpha

So be warned: no tests yet, totally untested

# Usage

`nginx-reloader <POLL_EVERY_SECONDS> [WATCHED_DIR...] -- [NGINX_CLI_OPTION...]`

e.g. command:

`nginx-reloader 10 /etc/nginx/conf.d`

- starts nginx
- checksums the `*.conf` files in `/etc/nginx/conf.d` directory every `10` seconds
- reloads nginx on every change

signals forwarded to nginx:

SIGHUP,
SIGINT,
SIGTERM,
SIGQUIT,
SIGKILL,
SIGABRT
