# HolyFlow

HolyFlow - это stateless сервис для организации хранения текстов христианских песен, их ритмичности и примеров mp3.

## Структура проекта

```
holyflow/
├── backend/          # Бэкенд приложение на Go
│   ├── cmd/          # Точка входа в приложение
│   ├── internal/     # Внутренние пакеты
│   │   ├── config/   # Конфигурация приложения
│   │   ├── database/ # Подключение к базе данных
│   │   ├── handlers/ # HTTP обработчики
│   │   ├── middleware/ # Middleware функции
│   │   ├── models/   # Модели данных
│   │   ├── server/   # HTTP сервер
│   │   └── services/ # Бизнес-логика
│   └── go.mod        # Зависимости Go
├── frontend/         # Фронтенд приложение на React
│   ├── public/       # Статические файлы
│   └── src/          # Исходный код React приложения
├── init-scripts/     # Скрипты инициализации
└── docker-compose.yml # Конфигурация Docker Compose
```

## Архитектура

### Backend
- Stateless Golang приложение
- Подключение к S3 для хранения mp3, обложек треков и аватарок профилей
- Логика личного кабинета для сбора любимых треков
- Аутентификация через логин/пароль и OIDC (GitLab)
- Подключение к внешнему SMTP серверу для восстановления пароля и подтверждения регистрации
- PostgreSQL в качестве базы данных

### Frontend
- Простое, легкое и лаконичное приложение
- Возможность добавления текстов песен, ритмичности и mp3 файлов для авторизованных пользователей
- Реализовано на React с использованием React Router для навигации
- Поддержка адаптивного дизайна для различных устройств

## Технологии

- **Backend**: Golang, Gin, GORM, PostgreSQL, MinIO
- **Frontend**: React, React Router, Axios
- **Аутентификация**: JWT, OAuth2 (GitLab)
- **Хранение файлов**: MinIO (S3 совместимое хранилище)

## Запуск приложения

### Требования
- Docker
- Docker Compose

### Запуск

```bash
# Клонирование репозитория
git clone <repository-url>
cd holyflow

# Запуск всех сервисов
docker-compose up -d

# Приложение будет доступно по адресу:
# Frontend: http://localhost:3000
```

### Переменные окружения

Backend поддерживает следующие переменные окружения:

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `DB_HOST` | Хост базы данных | `localhost` |
| `DB_PORT` | Порт базы данных | `5432` |
| `DB_USER` | Пользователь базы данных | `holyflow` |
| `DB_PASSWORD` | Пароль базы данных | `wkwz2i0uqceuudguf4kh` |
| `DB_NAME` | Имя базы данных | `holyflow` |
| `S3_ENDPOINT` | Endpoint S3 | `http://localhost:9000` |
| `S3_ACCESS_KEY` | Access key S3 | `minioadmin` |
| `S3_SECRET_KEY` | Secret key S3 | `014lq0tfcsm5s67gkwi1` |
| `S3_BUCKET` | Bucket S3 | `holyflow` |
| `S3_REGION` | Регион S3 | `us-east-1` |
| `JWT_SECRET` | Секретный ключ для JWT | `7f3c9a2e1b6d4f8a0c5e9d3b7a1f6c2e8d4b0a9f5c3e7d1b6a8f2c4e9d0b7` |
| `SMTP_HOST` | SMTP сервер | `` |
| `SMTP_PORT` | Порт SMTP | `587` |
| `SMTP_USERNAME` | Пользователь SMTP | `` |
| `SMTP_PASSWORD` | Пароль SMTP | `` |
| `SMTP_FROM` | Адрес отправителя | `noreply@holyflow.com` |
| `GITLAB_CLIENT_ID` | Client ID для GitLab OIDC | `` |
| `GITLAB_CLIENT_SECRET` | Client Secret для GitLab OIDC | `` |
| `GITLAB_REDIRECT_URL` | URL перенаправления для GitLab OIDC | `http://localhost:8080/api/v1/auth/gitlab/callback` |

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация нового пользователя
- `POST /api/v1/auth/login` - Вход по логину/паролю
- `POST /api/v1/auth/gitlab/login` - Начало OAuth через GitLab
- `GET /api/v1/auth/gitlab/callback` - Callback от GitLab
- `POST /api/v1/auth/forgot-password` - Восстановление пароля
- `POST /api/v1/auth/reset-password` - Сброс пароля

### Пользователи
- `GET /api/v1/users/me` - Получение профиля текущего пользователя
- `PUT /api/v1/users/me` - Обновление профиля
- `GET /api/v1/users/favorites` - Получение списка избранных песен
- `POST /api/v1/users/favorites/{songId}` - Добавление песни в избранное
- `DELETE /api/v1/users/favorites/{songId}` - Удаление песни из избранного

### Песни
- `GET /api/v1/songs` - Получение списка песен
- `GET /api/v1/songs/{id}` - Получение информации о песне
- `POST /api/v1/songs` - Создание новой песни (требуется аутентификация)
- `PUT /api/v1/songs/{id}` - Обновление информации о песне (требуется аутентификация)
- `DELETE /api/v1/songs/{id}` - Удаление песни (требуется аутентификация)
- `POST /api/v1/songs/{id}/upload-mp3` - Загрузка MP3 файла (требуется аутентификация)
- `POST /api/v1/songs/{id}/upload-cover` - Загрузка обложки (требуется аутентификация)

## Разработка

### Backend

```bash
cd backend

# Установка зависимостей
go mod tidy

# Запуск приложения
go run cmd/main.go
```

### Frontend

```bash
cd frontend

# Установка зависимостей
npm install

# Запуск приложения в режиме разработки
npm start

# Приложение будет доступно по адресу:
# http://localhost:3000
```

### Сборка Frontend

```bash
# Сборка для продакшена
npm run build
```

## Лицензия

MIT
