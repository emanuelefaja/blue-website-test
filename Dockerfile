FROM golang:1.24.4
ADD . /code
WORKDIR /code
RUN go mod download
CMD ["go", "run", "main.go"]