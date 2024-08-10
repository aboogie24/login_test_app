FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY main.go ./
COPY ./migrations/ ./migrations/
COPY .env.development.local ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /login_test_app

EXPOSE 8080

CMD ["/login_test_app"]