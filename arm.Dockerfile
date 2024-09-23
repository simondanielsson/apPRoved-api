FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/api cmd/main.go

EXPOSE 8090

CMD ["./bin/api"]
