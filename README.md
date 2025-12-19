# OAuth 2.0 Authorization Service

Сервис авторизации, реализующий OAuth 2.0-подобный flow с использованием **JWT access-токенов** и **refresh-токенов с ротацией**.  
Предназначен для использования как централизованный auth-сервис в микросервисной архитектуре.

---

## Возможности

- JWT Access Token
- Refresh Token с **ротацией** (one-time use)
- Хранение refresh-токенов в БД
- Защита от повторного использования refresh-токенов
- HttpOnly cookies для refresh-токенов
- Stateless access-токены
- Запуск через **Docker Compose**

---

## Запуск через Docker Compose

```bash
docker compose up --build