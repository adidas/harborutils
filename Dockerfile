FROM golang:1.15.8

WORKDIR /build/harborutils

COPY . ./
RUN go mod download

RUN go build  -o harborutils -ldflags "-linkmode external -extldflags -static" .

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /build/harborutils/harborutils /harborutils
ENTRYPOINT ["/harborutils"]
