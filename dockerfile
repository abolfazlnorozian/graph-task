    FROM golang:1.22-alpine AS builder

    WORKDIR /app
    
    # Dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -o server ./cmd/server
    
    
    # ---------- Runtime stage ----------
    FROM gcr.io/distroless/base-debian12
    
    WORKDIR /app
    
    COPY --from=builder /app/server /app/server
    
    EXPOSE 8080
    
    USER nonroot:nonroot
    
    ENTRYPOINT ["/app/server"]
    