FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/samkreter/DSA-Workshop/cmd/unzip/ 
COPY . /go/src/github.com/samkreter/DSA-Workshop/cmd/unzip/
RUN go get github.com/hubertsui/better-etling/util/extract
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run .

FROM golang:1.9.2
WORKDIR /root/
COPY --from=builder /go/src/github.com/samkreter/DSA-Workshop/cmd/unzip/run .
CMD ["./run"]