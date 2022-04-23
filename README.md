# price_hunter_ITMO
Backend for mobile app Price Hunter, final qualifying work.

It provides API for mobile app
API methods listed below (not exactly accurate, **WIP**)

[//]: # (TODO: add methods, paths, variables and error types from GoogleTable)

### Регистрация
POST-запрос с полями:
{
“email”: *,
“username”: *,
“password”: *
}

Ответ сервера с кодом
HTTP 200 с полями:
{
“token”: *
}

### Авторизация
POST-запрос с полями: {
“login”: *,
“password”: *
}

Ответ сервера с кодом
HTTP 200 с полями:
{
“token”: *
}

### Выход из аккаунта
POST-запрос с полями: {
“authorization”: *
}

Ответ сервера с кодом
HTTP 200

### Поиск игр
GET-запрос с аргументами:
“string”, “sortby”, “tag”

Ответ сервера с кодом
HTTP 200 полями:
{
“games”: [
“cover”: *,
“name”: *,
“publisher”: *,
“tags”: [*],
“id”: *
]
}

### Изменение адреса электронной почты
POST-запрос с полями: {
“new_email”: *
}

Ответ сервера с кодом
HTTP 200

### Изменение пароля
POST-запрос с полями: {
“current_password”: *,
“new_password”: *
}

Ответ сервера с кодом
HTTP 200  с полями:
{
“token”: *
}

### Получение информации об отдельной игре
GET-запрос с полями: {
“id”: *
}

Ответ сервера с кодом
HTTP 200 полями:
{
“cover”: *,
“name”: *,
“publisher”: *,
“description”: *,
“is_favorite”: *,
“tags”: [*],
“id”: *
}

### Добавление игры в избранное
GET-запрос с полями: {
“id”: *
} (?)

Ответ сервера с кодом
HTTP 200

### Получение всех избранных игр
GET-запрос (?)

Ответ сервера с кодом
HTTP 200 полями:
{
“games”: [
“cover”: *,
“name”: *,
“publisher”: *,
“tags”: [*],
“id”: *
]
}
