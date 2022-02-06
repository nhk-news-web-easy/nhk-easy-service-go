FROM golang:1.17 AS build
WORKDIR /app
COPY . /app
RUN go build -o nhk

FROM ubuntu
COPY --from=build /app/nhk .

EXPOSE 8080

CMD ["./nhk"]