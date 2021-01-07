FROM golang:1.13-alpine3.12
WORKDIR /workspace
COPY . /workspace

# install essential tools
RUN apk add curl bash ncurses jq bc make git expect

# install go command
RUN make all

# remove redundant things
RUN rm -r /workspace/*
