FROM golang:1.13 as builder

WORKDIR /go/src/logging-operator

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /go/src/logging-operator

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go/src/logging-operator/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
