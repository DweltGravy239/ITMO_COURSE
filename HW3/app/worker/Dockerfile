FROM golang:1.22 AS build
WORKDIR /worker
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o worker .
FROM scratch
COPY --from=build /worker/worker .
EXPOSE 8080
CMD ["/worker"]
