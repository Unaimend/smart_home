# Use the official Go image for Go 1.21
FROM golang:1.23

# Set environment variables for Go
ENV GO111MODULE=on \
    GOPATH=/go \
    PATH=$GOPATH/bin:/usr/local/go/bin:$PATH 
# Set the working directory inside the container
WORKDIR /app

# Copy the current project files to the container
COPY ./middleware/* /app
COPY ./middleware/cert/cert.pem /app
COPY ./middleware/cert/key.pem /app


# Run any required initialization or dependencies installation (optional)
RUN go mod tidy
RUN go mod download
# Define the default command to run
CMD ["go", "run", "main.go"]

