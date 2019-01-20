FROM golang:1.11-alpine3.8

RUN apk update && \
    apk add mongodb git build-base

RUN go get -u github.com/kardianos/govendor

RUN mkdir -p /data/db

WORKDIR /go/src/app
COPY . .

RUN govendor install +vendor,^program
RUN go get -d -t -v ./...
RUN go install -v ./...

WORKDIR /
COPY ./entrypoint.sh .

EXPOSE 80 27017

ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "run app" ]
