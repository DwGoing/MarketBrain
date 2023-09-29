FROM golang:1.20

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

COPY ./src  /
WORKDIR /
RUN go build -o MarketBrain main.go

ENTRYPOINT ["./MarketBrain", "service", "-t", "FUNDS"]

EXPOSE 8080 8081