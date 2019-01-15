FROM golang:1.11-alpine3.8

RUN apk update && \
    apk add mongodb git build-base

RUN mkdir -p /data/db

WORKDIR /go/src/app
COPY ./src .

COPY ./test_accounts_291218/data/data.zip /tmp/data/data.zip

RUN go get -d -v ./...
RUN go get -t -v ./...
RUN go install -v ./...

WORKDIR /
COPY ./entrypoint.sh .

ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "run app" ]
