ARG GOLANG_VERSION=1.11-alpine
FROM golang:${GOLANG_VERSION} as builder
LABEL maintainer "pgillich ta gmail.com"

ARG REPO_NAME="pgillich/airport-distance"
ARG BIN_PATH="/airport-distance"

COPY . /src
WORKDIR /src
RUN apk add git && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ${BIN_PATH} .

# Making minimal image (only one binary)
# TODO make CA certs manually, instead of using alpine
#FROM scratch
FROM alpine

ARG BIN_PATH="/airport-distance"
ARG RECEIVE_PORT="8080"

COPY --from=builder ${BIN_PATH} "/airport-distance"
ENTRYPOINT ["/airport-distance", "-service"]

RUN apk add ca-certificates

EXPOSE ${RECEIVE_PORT}