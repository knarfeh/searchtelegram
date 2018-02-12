FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/knarfeh/searchtelegram/
COPY . /go/src/github.com/knarfeh/searchtelegram/

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 8.9.1
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

EXPOSE 80 5000

CMD ["searchtelegram", "run"]
