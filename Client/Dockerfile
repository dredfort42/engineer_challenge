FROM golang:latest AS env

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app
RUN go mod download

FROM env
COPY . /app
RUN go build -o /app/client /app/cmd/client/main.go

CMD ["/app/client"]