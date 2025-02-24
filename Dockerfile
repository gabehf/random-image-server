## syntax=docker/dockerfile:1
FROM golang:1.23
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN mkdir -p /images
RUN CGO_ENABLED=0 GOOS=linux go build -o /random-image-server
CMD ["/random-image-server"]