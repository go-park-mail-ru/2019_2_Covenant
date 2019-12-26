FROM golang
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends postgresql-client
COPY . .
RUN go mod download
RUN make build
EXPOSE 8000
CMD sh until pg_isready --username=postgres --host=postgres; do sleep 1; done \
    && psql --username=postgres --host=postgres --list
ENTRYPOINT ./bin/api

