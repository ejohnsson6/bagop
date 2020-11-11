# build stage
FROM golang:latest AS build-env
ADD . /src
RUN cd /src && go build -o bagop

# final stage
FROM ubuntu:latest
WORKDIR /app
COPY --from=build-env /src/bagop /app/
COPY cron.sh /app/

RUN apt-get update && apt-get -y install cron

CMD ./cron.sh