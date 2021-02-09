# build stage
FROM golang:latest AS build-env
ADD . /src
RUN cd /src && go build -o bagop github.com/swexbe/bagop/cmd/bagop

# final stage
FROM ubuntu:latest
WORKDIR /app
COPY --from=build-env /src/bagop /app/
COPY cron.sh /app/

RUN apt-get update 
RUN apt-get -y install cron
RUN apt-get -y install ca-certificates

CMD ./cron.sh