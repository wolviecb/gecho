FROM golang:1.19 as build

ENV CGO_ENABLED=0
WORKDIR /gecho

COPY go.* main.go /gecho/
RUN go get ./... && \
    go build

FROM scratch

COPY --from=build /gecho/gecho /

ENTRYPOINT [ "/gecho" ]