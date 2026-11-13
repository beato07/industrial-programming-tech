# Notes REST API (Практическая работа №8)

RESTful веб-сервис для управления заметками, разработанный на языке **Go** с использованием микрофреймворка **Chi** и базы данных **MongoDB**.

## Требования

Для запуска и тестирования проекта необходимы:
* Go >= 1.21
* Docker и Docker Compose
* PowerShell (`Invoke-RestMethod`), `curl` или Postman

## Команды запуска

```bash
# 1. Запуск MongoDB в Docker
docker compose up -d

# 2. Создание файла конфигурации среды
cp .env.example .env  # PowerShell: Copy-Item .env.example .env

# 3. Запуск API-сервера
go run ./cmd/api```


## Примеры запросов

```bash
# 1. Создание заметки (POST)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notes" -Method POST -ContentType "application/json" -Body '{"title":"Go Programming","content":"Learning REST API with Gin/Chi and MongoDB"}'

# 2. Получение списка заметок с пагинацией (GET)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notes?limit=10&skip=0" -Method GET

# 3. Поиск заметок по тексту (GET с параметром q)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notes?q=Database" -Method GET

# 4. Получение статистики по базе данных (Aggregation Pipeline)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notes/stats" -Method GET

# 5. Получение заметки по ID (GET)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notes/stats" -Method GET

# 6. Частичное обновление заметки (PATCH)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notes/6af7007adf716b844900d6d" -Method PATCH -ContentType "application/json" -Body '{"title":"Updated Go Programming Title"}'

# 7. Удаление заметки (DELETE)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notes/6af7007adf716b844900d6d" -Method DELETE

# Запуск тестов
go test -v ./...