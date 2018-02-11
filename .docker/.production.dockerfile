FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/knarfeh/searchtelegram/
COPY . /go/src/github.com/knarfeh/searchtelegram/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o searchtelegram .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/bitbucket.org/ee-book/vhfw/vhfw /bin/
CMD ["vhfw", "serve"]
# ENTRYPOINT "/bin/vhfw serve"
