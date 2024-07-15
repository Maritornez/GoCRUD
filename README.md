# WebAPI на языке Go

Проект представляет из себя два Docker-контейнера. Один из них содержит базу данных Reindexer, а другой - WebAPI. Они запускаются с помощью Docker Compose. Есть три сущности: company, man и tip. Tip относится к какому-то man, а man в свою очередь относится к какому-то company. Реализованы операции CRUD над этими документами. У метода GET есть пагинация. Используется фреймворк Gin, потому что он поддерживает промежуточное ПО, с ним легко работать со сбоями, легко реализовать авторизацию. Кэширование документов реализовано посредством BigCache. Использовать Redis и развертвать его в отдельном контейнере не увидел смысла, потому что что Redis, что используемый Reindexer используют для хранения данных оперативную память, поэтому будут иметь примерно одинаковое время отклика. Ради того, чтобы добиться уменьшения времени отклика, данные хранятся в памяти самого приложения.

![изображение](https://github.com/Maritornez/Golang_CRUD/assets/62441435/b6f11d0b-837f-4483-9f39-a66587ea395c)


# Как запустить приложение

Контейнеры загружены на DockerHub.

Запуск контейнера с базой данных Reindexer
```
docker run -d -p 6534:6534 -p 9088:9088 --net involta --name reindexer reindexer/reindexer
```
Запуск контейнера с серверным приложением
```
docker pull ejrglkenr/go_crud_v2
docker run -d -p 8080:8080 --net involta --name go_crud_v2 ejrglkenr/go_crud_v2
```

# Как взаимодействовать с приложением

Endpoints для тестирования с помощью Postman:

Company:
- POST   http://localhost:8080/company?limit=10&offset=0
- GET    http://localhost:8080/companies
- GET    http://localhost:8080/company/{id}
- PATCH  http://localhost:8080/company/{id}
- DELETE http://localhost:8080/company/{id}

Man:
- POST   http://localhost:8080/man?limit=10&offset=0
- GET    http://localhost:8080/men
- GET    http://localhost:8080/man/{id}
- PATCH  http://localhost:8080/man/{id}
- DELETE http://localhost:8080/man/{id}

Tip:
- POST   http://localhost:8080/tip?limit=10&offset=0
- GET    http://localhost:8080/tips
- GET    http://localhost:8080/tip/{id}
- PATCH  http://localhost:8080/tip/{id}
- DELETE http://localhost:8080/tip/{id}

*Вместо id нужно подставить идентификатор документа*


Схемы документов (для того, чтобы по ним создавать свои json для HTTP-запросов) представлена в данном репозитории в директории `assets`
