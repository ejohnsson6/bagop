# build stage
FROM golang:1.17-alpine AS build-env
ADD . /src
RUN apk update && apk add --no-cache git
RUN cd /src && go build -o bagop github.com/swexbe/bagop/cmd/bagop

# final stage
FROM alpine
WORKDIR /home/root
COPY --from=build-env /src/bagop /app/
# Add certs
RUN apk add -U --no-cache ca-certificates
# If you need to use cli
RUN apk add bash
RUN ["ln", "-s", "/app/bagop", "/usr/bin/bagop"]
VOLUME /var/bagop

CMD /app/bagop -s