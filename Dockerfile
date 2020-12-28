# Start from the latest golang base image
FROM golang:1.15.6-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . /app

RUN cd /app && go install 
