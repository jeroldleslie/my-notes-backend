# Compile app binary
FROM golang:latest as build-env

WORKDIR /go/src/github.com/jeroldleslie/my-notes-backend
ARG GO_APP_LOCATION
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal internal
COPY cmd/${GO_APP_LOCATION} cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix cgo -o ~/app ./cmd/*.go

# Run app in scratch
FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /root/app /

EXPOSE 8000
CMD ["/app"]