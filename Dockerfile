FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o fp_pinjaman_online

EXPOSE 8080

ENTRYPOINT ["/app/fp_pinjaman_online"]