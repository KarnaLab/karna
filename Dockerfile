# Start from the latest golang base image
FROM golang:1.15.6-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app
COPY . /app
RUN cd /app && go build 

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app .
RUN ls
CMD ["./karna"] 