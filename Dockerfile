FROM golang:1.23.4
LABEL authors="itsme"

#set the current working directory inside the container
WORKDIR /app

#Copy go.mod and go .sum files
COPY go.mod  go.sum ./

#Download all dependencies. Dependencies will be caches if the go.mod and go.sum files are not changed

RUN go mod download

# COPY the source code into the container

COPY . .

#Build the Go app
RUN go build -o main .

CMD ["./main"]