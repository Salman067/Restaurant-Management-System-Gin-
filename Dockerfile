# syntax=docker/dockerfile:1
# Stage 1: Build stage
FROM golang:1.19-alpine AS builder

# Install the 'upx' package
RUN apk add upx
# testing hook
# Set the working directory
WORKDIR /app

# Create a user and group for running the application
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the Go application with specified flags and compress it with 'upx'
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod GOOS=linux go build -ldflags="-w -s" -a -o /docker-gs-ping .
RUN upx --best --lzma /docker-gs-ping

# Stage 2: Final stage
FROM alpine AS final

# Import the user and group files from the builder stage
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the builder stage
COPY --from=builder /app/base.env /base.env
COPY --chown=nobody:nobody --from=builder /docker-gs-ping /docker-gs-ping
COPY --from=builder /app/localize /localize

# Create and set permissions for the attachments directory
#RUN mkdir /attachments; \
#    chown -R nobody:nobody /attachments; \
#    chmod -R 777 /attachments

# Specify a volume for attachments
VOLUME ["attachments"]

# Perform any further action as an unprivileged user
USER nobody:nobody

# Expose the container's port
EXPOSE 8089

# Run the compiled binary as the entrypoint
ENTRYPOINT ["/docker-gs-ping", "server"]

