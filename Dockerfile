FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main
RUN ls
EXPOSE 3000

CMD ["/main"]
