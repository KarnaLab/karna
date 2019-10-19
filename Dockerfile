# Start from the latest golang base image
FROM golang:1.13.2-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

RUN git clone https://github.com/karbonn/karna.git

COPY . /app

RUN cd /app/karna && go install 
