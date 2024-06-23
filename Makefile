.PHONY: all build clean

all: server

server:
	# Запуск https сервера на порту 3333
	 go run cmd/ex00-server/main.go --tls-port 3333