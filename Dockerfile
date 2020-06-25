FROM golang
WORKDIR /myapp
COPY go.mod .
COPY wserver.go .
RUN go mod download && go build 
EXPOSE 5000
ENTRYPOINT ["./server"]

