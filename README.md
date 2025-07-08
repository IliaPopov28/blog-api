# Blog API

Простое REST API для блога на Go (Gin + GORM + PostgreSQL)

## Возможности
- Регистрация и аутентификация пользователей (JWT)
- CRUD для постов (создание, чтение, обновление, удаление)
- Защита эндпоинтов через JWT

## Стек
- Go
- Gin
- GORM
- PostgreSQL

## Быстрый старт

### 1. Клонируй репозиторий
```sh
git clone <your-repo-url>
cd blog-api
```

### 2. Установи зависимости
```sh
go mod tidy
```

### 3. Создай файл `.env` с переменными окружения
```
DATABASE_DSN=host=localhost user=postgres password=yourpassword dbname=blog_api_db port=5432 sslmode=disable
JWT_SECRET=your_secret_key
```

### 4. Запусти приложение
```sh
go run main.go
```

## Примеры запросов

### Регистрация
```sh
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"myuser", "password":"mypassword"}'
```

### Логин
```sh
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"myuser", "password":"mypassword"}'
```

### Создать пост
```sh
curl -X POST http://localhost:8080/posts \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Заголовок", "content":"Текст поста"}'
```

### Получить все посты
```sh
curl -H "Authorization: Bearer <your_token>" http://localhost:8080/posts
```

## Переменные окружения
- `DATABASE_DSN` — строка подключения к PostgreSQL
- `JWT_SECRET` — секрет для подписи JWT

