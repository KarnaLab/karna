# Start from the latest golang base image
FROM golang:1.13.2-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . /app

RUN cd /karna && go install 
