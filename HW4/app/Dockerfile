FROM golang:1.22 AS build
WORKDIR /batch
COPY go.mod .
COPY main.go .
RUN CGO_ENABLED=0 go build -o batch .
FROM scratch
COPY --from=build /batch/batch .
USER 10001
CMD ["/batch"]
