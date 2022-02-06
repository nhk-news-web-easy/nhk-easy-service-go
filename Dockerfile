FROM golang:1.17

WORKDIR /app
COPY . /app

RUN go build -o nhk

EXPOSE 8080

CMD ["./nhk"]