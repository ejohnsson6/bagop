# build stage
FROM golang:1.15-alpine AS build-env
ADD . /src
RUN apk update && apk add --no-cache git
RUN cd /src && go build -o bagop github.com/swexbe/bagop/cmd/bagop

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/bagop /app/
# Add certs
RUN apk add -U --no-cache ca-certificates
COPY run_and_wait.sh /app/

CMD ./run_and_wait.sh