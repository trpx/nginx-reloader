

FROM golang:1.12-alpine

WORKDIR /tmp

RUN apk add git --no-cache \
    && go get -v github.com/trpx/nginx-reloader


FROM nginx:1.16-alpine

COPY --from=0 /go/bin/nginx-reloader /usr/local/bin/

ENTRYPOINT ["nginx-reloader", "1", "/etc/nginx/conf.d", "--", "-g", "daemon", "off"]