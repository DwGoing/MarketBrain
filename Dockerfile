FROM golang:1.20-alpine

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

COPY ./bin  /
WORKDIR /

RUN MarketBrain service -t FUNDS

EXPOSE 8080 8081