FROM golang:1.20 as build

ARG TARGETOS
ARG TARGETARCH

# Github Actions build labels
ARG BUILD_DATE

ENV BUILD_DATE=$BUILD_DATE
ENV GITHUB_SHA=$GITHUB_SHA

WORKDIR /gecho

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY gecho.go gecho.go
COPY tools.go tools.go
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o gecho *.go


FROM scratch

COPY --from=build /gecho/gecho /

ENTRYPOINT [ "/gecho" ]
