# Exchanger

## Описание проекта

Exchanger — микросервис для управления данными криптовалют. Проект реализует GRPC сервер, взаимодействующий с PostgreSQL для хранения данных, и поддерживает трейсинг с помощью Jaeger.

## Основные возможности

- GRPC API для управления данными.
- Логирование с использованием Zap.
- Трейсинг запросов через OpenTelemetry и Jaeger.
- Автоматическая миграция базы данных (GORM).

## Установка и запуск

### Зависимости

Перед началом работы установите:

- [Go](https://go.dev/) >= 1.20
- [Docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/)
- [PostgreSQL](https://www.postgresql.org/)

### Сборка и запуск

1. **Склонируйте репозиторий:**
   ```bash
   git clone https://studentgit.kata.academy/gk/exchanger.git
   cd exchanger
   ```

2. **Соберите бинарный файл:**
   ```bash
   make build
   ```

3. **Запустите тесты:**
   ```bash
   make test
   ```

4. **Соберите Docker-образ:**
   ```bash
   make docker-build
   ```

5. **Запустите приложение:**
   ```bash
   make run
   ```

6. **Проверьте код линтером:**
   ```bash
   make lint
   ```

## Конфигурация

Основные параметры задаются через конфигурационный файл:

- Настройки PostgreSQL: `DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_PORT`.
- URL для Jaeger: `http://jaeger:14268/api/traces`.
- Порт GRPC: `50051`.

## Разработка

Для управления проектом используйте команды Makefile:

- `make build` — сборка бинарного файла.
- `make test` — запуск тестов.
- `make docker-build` — сборка Docker-образа.
- `make run` — запуск через Docker Compose.
- `make lint` — проверка кода линтером.

## Технологии

- **Go**: язык программирования.
- **PostgreSQL**: база данных.
- **GRPC**: интерфейс API.
- **OpenTelemetry**: инструмент трейсинга.
- **Jaeger**: система распределённого трейсинга.
- **Zap**: библиотека логирования.

