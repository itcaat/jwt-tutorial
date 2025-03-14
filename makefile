# Variables
PATH_NGINX_INGRESS = nginx-ingress
PATH_JWT_APP = jwt-app
PATH_UI = ui

# Argument parsing
args = $(foreach a,$($(subst -,_,$1)_args),$(if $(value $a),$a="$($a)"))

init:
	cd $(PATH_NGINX_INGRESS); mkcert -key-file key.pem -cert-file cert.pem 'localhost.devopsbrain.ru' '*.localhost.devopsbrain.ru' ;

run-compose:
	docker compose up --build;

run-compose-daemon:
	docker compose up --build -d;

run-api:
	cd $(PATH_JWT_APP); go run cmd/main.go;

test:
	cd $(PATH_JWT_APP); go test ./internal/tests;

ls:
	find . -type f -not -path "./.git/*" -exec echo "===== {} =====" \; -exec cat {} \;\n
