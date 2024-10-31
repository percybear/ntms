# Dockerfile using multistage build to create a small image.
# Build with golang:1.23-alpine as it contains the golang compiler and its dependencies.
# These take up disk space, and we don't need them in the final image.
FROM golang:1.23-alpine AS build
WORKDIR /go/src/ntms
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/ntms ./cmd/ntms


# Use scratch as the final image, a minimal image that contains only the runtime environment.
# FROM scratch
FROM busybox:latest
COPY --from=build /go/src/ntms/config /etc/options/ntms
COPY --from=build /go/bin/ntms /bin/ntms
# ENTRYPOINT ["httpd", "-f", "-p", "8080"]
ENTRYPOINT ["ntms"]
