FROM golang:1.22
WORKDIR /app
COPY main.go .
COPY go.mod .
RUN CGO_ENABLED=0 go build -o api .
EXPOSE 8080
CMD ["/app/api"] 
