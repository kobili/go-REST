FROM golang:1.23.3-alpine

# establish the working directory in the container
WORKDIR /app

# download air
RUN go install github.com/cosmtrek/air@v1.51.0

# download delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@v1.23.1

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY .air.toml .

# Entry point
CMD ["air", "."]
