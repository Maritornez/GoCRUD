# WebAPI на языке Go

Проект представляет из себя два Docker-контейнера. Один из них содержит базу данных Reindexer, а другой - WebAPI.

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
- GET    http://localhost:8080/men
- GET    http://localhost:8080/man/{id}
- POST   http://localhost:8080/man?limit=10&offset=0
- PATCH  http://localhost:8080/man/{id}
- DELETE http://localhost:8080/man/{id}

*Вместо id нужно подставить идентификатор документа*


Схема документа (для того, чтобы по ней создавать свои документы для HTTP-запросов) представлена в данном репозитории в директории assets/TestInputMan.json
