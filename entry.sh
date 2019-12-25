until pg_isready --dbname=postgres --username=postgres --host=postgres; do sleep 1; done \
        && psql --username=postgres --host=postgres --list
./bin/api          
