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
- GET    http://localhost:8080/mans
- GET    http://localhost:8080/mans/id (вместо id нужно подставить идентификатор документа)
- POST   http://localhost:8080/mans
- PATCH  http://localhost:8080/mans/id (вместо id нужно подставить идентификатор документа)
- DELETE http://localhost:8080/mans/id (вместо id нужно подставить идентификатор документа)

Схема документа представлена в данном репозитории в директории assets/TestInputMan.json
