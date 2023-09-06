FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -o /cloud-drive

EXPOSE 8080

CMD ["/cloud-drive"]