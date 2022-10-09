FROM golang:1.16-alpine

WORKDIR /app

COPY . ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

RUN go build -o /picture-oss-proxy

EXPOSE 3000

CMD [ "/picture-oss-proxy" ]