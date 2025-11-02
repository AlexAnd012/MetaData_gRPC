# MediaMeta (gRPC, Go, Postgres)

Мини-сервис для учёта метаданных файлов.
Xраним только сведения о них (имя, размер,  владелец, дата создания).
Основной демо-поток: клиент передаёт локальный путь к файлу, сервер сам читает метаданные и сохраняет их в Postgres.

## Возможности

**gRPC API:**

**HealthCheck**

**CreateMetadataFromPath** — создать запись по локальному пути (сервер читает os.Stat, MIME по расширению)

**GetMetadata**

**ListMetadata** (пагинация page_size + page_token как offset)

**UpdateMetadata** (частичное: filename, content_type)

**DeleteMetadata** (жёсткое удаление)

Простая архитектура: сервис ↔ репозиторий (Postgres).

**Миграции SQL** (упрощённая схема).

ручной метод CreateMetadata можно оставить или удалить.

Архитектура
.
├─ cmd/  
│  ├─ server/      # запуск gRPC-сервера  
│  └─ client/      # демонстрационный клиент (передаёт локальный путь)  
├─ internal/  
│  ├─ service/     # бизнес-логика gRPC (валидация, коды ошибок)  
│  └─ storage/     # интерфейс Repository и реализация Postgres  
├─ proto/          # .proto (контракт gRPC)  
├─ gen/go/         # сгенерированный Go-код из .proto  
└─ migrations/     # SQL-миграции  

## Стек ##

Go 1.24

protoc + плагины:

protoc-gen-go

protoc-gen-go-grpc

Docker (+ Docker Compose) — для Postgres

## Быстрый старт ##
**1) Поднять Postgres в Docker**
docker compose up -d db
docker compose ps   

**2) Применить миграции**
$cid = docker compose ps -q db
Get-Content -Raw .\migrations\0001_init.sql | docker exec -i $cid psql -U mm -d mediameta

## Схема ##

CREATE TABLE IF NOT EXISTS metadata (
  id           uuid    PRIMARY KEY,
  filename     text    NOT NULL,
  size_bytes   bigint  NOT NULL DEFAULT 0,
  content_type text,
  owner_id     text,
  created_at   bigint  NOT NULL
);

## 3) Сгенерировать gRPC-код ##
protoc -I proto `
  --go_out=gen/go --go_opt=paths=source_relative `
  --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative `
  proto/mediameta/v1/metadata.proto

## 4) Запустить сервер ##
$env:DB_DSN="postgres://mm:mm@localhost:5434/mediameta?sslmode=disable"
go run .\cmd\server
