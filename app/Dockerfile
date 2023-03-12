FROM golang:1.18

WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any dependencies required by the application
RUN go get -d -v ./...

# Build the application inside the container
RUN go build -o main .

# Expose the port on which the application will listen
EXPOSE 8080

# Define the command to run the application when the container starts
CMD ["/app/main"]
