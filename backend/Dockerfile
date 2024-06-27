FROM golang:1.22.4-alpine

WORKDIR /api

COPY ./main.go .

RUN go mod init api-tutorial && go mod tidy

EXPOSE 8080

CMD ["go", "run", "./main.go"]