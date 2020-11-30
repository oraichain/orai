FROM golang:1.13-alpine3.12
WORKDIR /workspace
COPY . /workspace
RUN init.sh