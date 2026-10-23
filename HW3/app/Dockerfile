FROM golang:1.22 AS build
WORKDIR /app
COPY api/go.mod api/
COPY api/main.go api/
COPY api/go.sum api/
COPY worker/go.mod worker/
COPY worker/main.go worker/
COPY worker/go.sum worker/
RUN CGO_ENABLED=0 go build -o api api/
RUN CGO_ENABLED=0 go build -o worker worker/
FROM scratch
COPY --from=build /app/api .
COPY --from=build /app/worker .
EXPOSE 8080
CMD ["/api", "/worker"]
