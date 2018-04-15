FROM python:3.6-rc-alpine
MAINTAINER He Jun knarfeh@outlook.com

# base pkgs
RUN apk --update add --no-cache openssl ca-certificates postgresql-dev 

# build pkgs
RUN apk add --no-cache gcc g++ python3-dev musl-dev make

# dev pkgs
RUN apk add curl

COPY e2e/requirements /requirements
RUN pip3 install -U pip \
    && pip install -r requirements/dev.txt
COPY ./e2e /src/
WORKDIR /src