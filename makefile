# Variables
PATH_JWT_APP = jwt-app
PATH_UI = ui

# Argument parsing
args = $(foreach a,$($(subst -,_,$1)_args),$(if $(value $a),$a="$($a)"))

init:
	cd $(PATH_UI); mkcert -key-file key.pem -cert-file cert.pem '127.0.0.1.nip.io' '*.127.0.0.1.nip.io' ;

run-compose:
	docker compose up --build;

run-compose-daemon:
	docker compose up --build -d;

run-api:
	cd $(PATH_JWT_APP); go run cmd/main.go;

test:
	cd $(PATH_JWT_APP); go test ./internal/tests;
