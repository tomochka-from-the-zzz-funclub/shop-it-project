# Makefile для сервиса search

# ────────────────────────────────────────
# Если есть .env — включаем его и экспортируем все переменные
# ────────────────────────────────────────
ifneq (,$(wildcard .env))
  include .env
  export $(shell sed -n 's/=.*//' .env)
endif

# ────────────────────────────────────────
# Переменные
# ────────────────────────────────────────
GO         := go
SWAG       := swag
LIQUIBASE  := liquibase
DOCKER_PATH   := docker-compose -f docker/docker-compose.yml
MAIN_PATH  := ./cmd/user_service/main.go
BINARY     := bin/search

.DEFAULT_GOAL := help

# ----------------------------------------
# Помощь
# ----------------------------------------
.PHONY: help
help: ## Показать справочное сообщение
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ----------------------------------------
# Установка зависимостей и инструментов
# ----------------------------------------
.PHONY: deps
deps: ## Установить Go-модули и необходимые CLI-инструменты, если их нет, и загрузить .env
	@echo "---- deps: Начало ----"
	@echo "→ Загрузка переменных из .env..."
	@if [ -f .env ]; then \
	  set -o allexport; \
	  . .env; \
	  set +o allexport; \
	  echo "   • .env загружен"; \
	else \
	  echo "   • .env не найден"; \
	fi
	@echo "→ Проверка Go-модулей..."
	@$(GO) mod tidy
	@$(GO) mod download
	@echo "→ Проверка инструментов..."
	@command -v $(SWAG) >/dev/null 2>&1 || { echo "  • Устанавливаю swag..."; $(GO) install github.com/swaggo/swag/cmd/swag@latest; }
	@echo "---- deps: Завершено ----"


# ----------------------------------------
# Сборка и запуск
# ----------------------------------------
.PHONY: build
build: ## Собрать бинарник сервиса
	@echo "---- build: Начало ----"
	@mkdir -p bin
	@$(GO) build -o $(BINARY) ./cmd/search_service
	@echo "---- build: Завершено ----"

.PHONY: run
run: build ## Собрать и запустить сервис (требует .env)
	@echo "---- run: Начало ----"
	@echo "Запуск сервиса..."
	@$(BINARY)
	@echo "---- run: Завершено ----"

.PHONY: start
start: build ## Собрать и запустить сервис в фоне
	@echo "---- start: Начало ----"
	@echo "Запуск сервиса в фоне..."
	@$(BINARY) & echo $$! > .search.pid
	@echo "PID сохранён в .search.pid: $$(cat .search.pid)"
	@echo "---- start: Завершено ----"

.PHONY: stop
stop: ## Остановить сервис, используя PID
	@echo "---- stop: Начало ----"
	@if [ -f .search.pid ]; then \
		kill `cat .search.pid` && echo "Сервис остановлен"; \
		rm .search.pid; \
	else \
		echo "PID-файл не найден, сервис может быть не запущен"; \
	fi
	@echo "---- stop: Завершено ----"

# ----------------------------------------
# Swagger-документация
# ----------------------------------------
.PHONY: swagger
swagger: ## Сгенерировать Swagger-документацию (требуется swag)
	@echo "---- swagger: Начало ----"
	@$(SWAG) init -g ./internal/app/search.go
	@echo "---- swagger: Завершено ----"

# ----------------------------------------
# Docker Compose
# ----------------------------------------
.PHONY: docker-up
docker-up: ## Запустить сервисы Docker Compose
	@echo "---- docker-up: Начало ----"
	@$(DOCKER_PATH) up -d
	@echo "---- docker-up: Завершено ----"

.PHONY: docker-down
docker-down: ## Остановить сервисы Docker Compose
	@echo "---- docker-down: Начало ----"
	@$(DOCKER_PATH) down
	@echo "---- docker-down: Завершено ----"

.PHONY: docker-down-v
docker-down-v: ## Остановить и удалить сервисы Docker Compose вместе с volume
	@echo "---- docker-down-v: Начало ----"
	@$(DOCKER_PATH) down -v
	@echo "---- docker-down-v: Завершено ----"

# ----------------------------------------
# Очистка
# ----------------------------------------
.PHONY: clean
clean: ## Удалить собранные бинарники
	@echo "---- clean: Начало ----"
	@rm -rf bin
	@echo "---- clean: Завершено ----"

# ----------------------------------------
# Полный цикл запуска
# ----------------------------------------
.PHONY: up-all
up-all: deps swagger docker-up run ## Установить зависимости, сгенерировать Swagger, поднять Docker и запустить сервис
	@echo "---- up-all: Завершено ----"

# ----------------------------------------
# Полная остановка всех сервисов
# ----------------------------------------
.PHONY: down-all
down-all: stop docker-down  ## Остановить сервис и Docker Compose
	@echo "---- down-all: Завершено ----"
