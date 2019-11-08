FROM golang
WORKDIR /app
COPY . .
RUN go mod download
RUN make build
EXPOSE 8000
ENTRYPOINT ["./bin/api"]
