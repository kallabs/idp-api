FROM golang:1.21.5-alpine3.18 as debug

RUN apk update && apk upgrade && \
    apk add --no-cache git dpkg gcc musl-dev

ENV GOPATH /go
ENV GO111MODULE on
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /usr/src
ADD ./src ./

RUN ls -la .

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go mod download
RUN go build -o /usr/app ./cmd/http.go

COPY ./dlv.sh /
RUN chmod +x /dlv.sh
ENTRYPOINT [ "/dlv.sh" ]

FROM alpine:3.18 as prod

RUN apk --update add postgresql-client

COPY --from=debug /usr/app /usr/app
COPY ./wait-for-db.sh /
RUN chmod +x /wait-for-db.sh

CMD /wait-for-db.sh /usr/app