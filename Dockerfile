# stage 1 build the app using alpine environment
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o main .

# stage 2 run the app using the file from stage 1
FROM alpine:3.23.2

WORKDIR /app

COPY --from=builder /app/main .

COPY /app ./app/

COPY .env .

EXPOSE 8080

CMD [ "./main" ]
