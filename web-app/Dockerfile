FROM golang:1.22.4-alpine

WORKDIR /web-app

COPY ./main.go .

RUN go mod init web-app && go mod tidy

EXPOSE 8080

CMD ["go", "run", "./main.go"]