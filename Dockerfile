FROM golang
RUN apt-get update && apt-get install -y postgresql-client
WORKDIR /app
COPY . .
RUN go mod download
RUN make build
EXPOSE 8000
CMD sh until pg_isready --username=postgres --host=postgres; do sleep 1; done \
	&& psql --username=postgres --host=postgres --list
ENTRYPOINT ["./bin/api"]
