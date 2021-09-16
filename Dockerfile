# Build the manager binary
FROM quay.io/cybozu/golang:1.17-focal as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY version.go version.go
COPY cmd/ cmd/
COPY hooks/ hooks/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pod-security-admission cmd/main.go

FROM scratch
WORKDIR /
COPY --from=builder /workspace/pod-security-admission .
USER 10000:10000

ENTRYPOINT ["/pod-security-admission"]
