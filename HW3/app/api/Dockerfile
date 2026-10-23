FROM golang:1.22 AS build
WORKDIR /api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o api .
FROM scratch
COPY --from=build /api/api .
EXPOSE 8080
CMD ["/api"]
