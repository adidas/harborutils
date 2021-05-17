FROM golang:1.15.8

RUN go get -u gorm.io/gorm gorm.io/driver/postgres github.com/spf13/cobra github.com/howeyc/gopass
WORKDIR /go/src/github.com/adidas/harborutils
COPY . .
WORKDIR /root
RUN go build -o /go/src/github.com/adidas/harborutils -ldflags "-linkmode external -extldflags -static" /go/src/github.com/adidas/harborutils 
RUN ls /root

FROM scratch
COPY --from=0 /go/src/github.com/adidas/harborutils/harborutils /harborutils
ENTRYPOINT ["/harborutils"]
