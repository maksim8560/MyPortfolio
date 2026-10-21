# Сборка и запуск портфолио: фронтенд встраивается в Go-бинарник.

VERSION := $(shell date +%Y%m%d-%H%M)

.PHONY: all web server run dev clean

all: web server

web:
	cd web && npm install --no-audit --no-fund && npm run build

server:
	cd server && go build -ldflags "-X main.version=$(VERSION)" -o portfolio .

run: all
	cd server && ./portfolio

dev:
	@echo "терминал 1: cd server && DEV_ORIGIN=http://localhost:5173 go run ."
	@echo "терминал 2: cd web && npm run dev"

clean:
	rm -rf server/dist server/portfolio
