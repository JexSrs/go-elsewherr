FROM golang:1.20-alpine

WORKDIR "/app"

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Compile
RUN mkdir bin
RUN go build -o bin ./...

# Run
CMD [ "./bin/src" ]