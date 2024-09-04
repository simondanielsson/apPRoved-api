FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download
COPY . ./

RUN make build-linux

EXPOSE 8090

CMD ["make", "run"]
