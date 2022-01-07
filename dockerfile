FROM golang:1.17.5-alpine3.15

WORKDIR /app

COPY . .

RUN go build -o wecom

EXPOSE 8080

CMD ./wecom