# rest-wallet

## Описание

Этот проект реализует REST API для работы с кошельком пользователя. Пользователи могут выполнять следующие операции:

* Пополнение кошелька
* Списание средств с кошелька
* Просмотр баланса кошелька

Используется Postgresql в качестве субд, Docker для контейнеризации, Redis для кэширования. Код покрыт info и debug логами.
Был сгенерирован swagger на реализованный API

## Требования

- Go 1.23
- PostgreSQL
- Docker
- Redis

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/MaximInnopolis/rest-wallet.git
cd rest-wallet
```

2. Соберите докер-билд:
```bash
make up-all
```

3. Проведите миграцию:
```bash
make migrate
```

API доступен по адресу <http:localhost:8080>

Swagger доступен по адресу <http://localhost:8080/docs/swagger/index.html>