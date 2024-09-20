ifneq (,$(wildcard ./app.env))
    include app.env
    export
endif

hello:
	echo "Hello, World"

run-migration:
	cd migrations && GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=$(POSTGRES_USER) port=$(POSTGRES_PORT) host=$(POSTGRES_HOST) dbname=$(POSTGRES_DB) password=$(POSTGRES_PASSWORD) sslmode=disable" goose up

down-migration:
	cd migrations && GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=$(POSTGRES_USER) port=$(POSTGRES_PORT) host=$(POSTGRES_HOST) dbname=$(POSTGRES_DB) password=$(POSTGRES_PASSWORD) sslmode=disable" goose down

docker:
    docker run --name user-service -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -p $(POSTGRES_PORT):$(POSTGRES_PORT) -d postgres 
