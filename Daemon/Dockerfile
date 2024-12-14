FROM golang:latest AS env

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app
RUN go mod download

FROM env
COPY . /app
RUN go build -o /app/daemon /app/cmd/daemon/main.go

CMD ["/app/daemon"]