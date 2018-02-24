FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/knarfeh/searchtelegram/
COPY . /go/src/github.com/knarfeh/searchtelegram/

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 8.9.1
ENV CGO_ENABLED 0
ENV GOOS linux
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.2/install.sh | bash \
    && source $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default \
    && npm install \
    && npm rebuild node-sass --force

# Set up our PATH correctly
ENV NODE_PATH $NVM_DIR/versions/node/v$NODE_VERSION/lib/node_modules
ENV PATH      $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

RUN make build

FROM alpine:latest

EXPOSE 80 5000

RUN apk update && \
  apk --no-cache add ca-certificates supervisor nginx

WORKDIR /root/
COPY --from=builder /go/bin/searchtelegram /bin/
COPY --from=builder /go/src/github.com/knarfeh/searchtelegram/*.sh /
COPY --from=builder /go/src/github.com/knarfeh/searchtelegram/conf/supervisord.conf /etc/supervisord.conf
COPY --from=builder /go/src/github.com/knarfeh/searchtelegram/conf/nginx.conf /etc/nginx/searchtelegram_nginx.conf
RUN mkdir -p /var/log/supervisor /var/log/searchtelegram /media/images /var/nginx/cache/aws
RUN chmod +x /*.sh

CMD ["/searchtelegram.sh"]
