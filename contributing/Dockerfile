FROM ubuntu:18.04

ARG NODE_VERSION=8.10.0
ARG GO_VERSION=1.18.1
ARG BUD_VERSION=main

# Install basic dependencies
RUN apt-get -qq update
RUN apt-get -qq -y install curl git make gcc g++

# Install Node.js
RUN curl https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz | tar -xz -C /usr/local --strip-components 1
RUN node -v

# Install Go
RUN curl -L https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz | tar -xz -C /usr/local --strip-components 1
RUN go version
ENV PATH "/root/go/bin:${PATH}"

# Install bud
RUN git clone https://github.com/livebud/bud /bud
WORKDIR /bud
RUN git checkout $BUD_VERSION
RUN make install
RUN go install .
RUN bud version

