# Сборка приложения
build:
	go build -o bin/main ./cmd

# Запуск unit-тестов
test:
	go test ./...

# Сборка Docker-образа с использованием docker-compose
docker-build:
	sudo docker-compose build

# Запуск приложения с использованием docker-compose
run:
	sudo docker-compose up

# Запуск линтера
lint:
	golangci-lint run



# Стандартная цель
.PHONY: build test docker-build run lint clean
