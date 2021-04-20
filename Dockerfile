FROM golang:latest

WORKDIR /go/src/httsproxy

COPY . .

EXPOSE 1080

CMD ["go", "run", "main.go", "-L=0.0.0.0:1080"]
