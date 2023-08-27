# Avito тестовое задание для DevOps (2023)

## Структура проекта

```
├── docker-compose.yml
├── app
│   ├── Dockerfile.goapp
│   └── main.go
├── nginx
│   ├── Dockerfile.nginx
│   └── nginx.conf
├── redis
│   ├── Dockerfile.redis
│   └── redis.conf
└──README.md
```
app/ - директория с Golang приложением \
nginx/ - директория с конфигурацией и Dockerfile для запуска Nginx в контейнере \
redis/ - директория с конфигурацией и Dockerfile для запуска REDIS в контейнере \
docker-compose - compose файл для запуска всех сервисов в контейнерах

В качестве базового образа используется debian:12.1-slim во всех сервисах \
Приложение запускается на порту 8089 \
Запросы проксируются с nginx на Golang приложение, которое делает непосредственно запросы к REDIS 

Приложение использует 3 эндпоинта: \
1. POST /set_key { "<key>": "<val>" } - добавление ключ-значение в REDIS
2. GET /get_key?key=<key> - добавление ключ-значение в REDIS
2. DELETE /set_key { "<key>": "<val>" } - удаление ключ-значение по ключу

## Описание конфигурации Golang приложения
В файле main.go присутствуют два параметра конфигурации: \
redisAddr - адрес REDIS с портом (Пример: 127.0.0.1:6379) \
listenAddrServ - порт на котором работает Golang приложение (Пример: :8000)
Значения в файле main.go по умолчанию:
```
var redisAddr string = "redis-go:6379"
var listenAddrServ string = ":8000"
```
redis-go - имя сервиса при запуске приложения с помощью docker-compose

## Описание Nginx конфигурации
Nginx проксирует все запросы по 3 endpoint к приложению Golang \
Обращение идет по DNS имени (имени сервиса) и объявлено в директиве upstream в nginx/nginx.conf \
По умолчанию используется имя сервиса goapp и следующая конфигурация: 
```
    upstream backend {
        server goapp:8000;
    }
```
В случае, если необходимо запустить приложение локально, либо будет использоваться другой сервер с приложением, необходимо заменить goapp на другой адрес

## Описание REDIS конфигурации
Обращения от Golang приложения идут к REDIS \
Основные параметры для изменения в redis/redis.conf
port - порт, на котором функционирует REDIS сервер (По умолчанию: 6379)

## Запуск проекта
Для запуска проекта локально необходимо:

1. Установленный Docker, docker-compose 
Если данное ПО не установлено, необходимо обратиться к официальной документации: https://docs.docker.com/engine/install/

2. Установленный Git. Если не установлен, обратиться к официальной документации: https://git-scm.com/downloads

3. Установленный curl для проверки корректности работы приложения (опционально)



Далее будет приведен пример запуска приложения на Linux, используемый дистрибутив: Ubuntu 20.04 \

Клонируем проект с Git и перейдем в директорию с проектом:
```
git clone https://github.com/baranovr605/avito-internship-devops.git
cd avito-internship-devops

```

Запустим проект средствами docker-compose в фоновом режиме:
```
docker-compose up -d
```

Если все успешно в консоли увидим следующие строки:
```
Creating redis-go ... done
Creating nginx-go ... done
Creating goapp    ... done
```

Можно провести проверку работы приложения. \
Ниже представлен пример с записью ключ-значение, получение данных по ключу и удаление ключа:
```
# curl -X POST localhost:8089/set_key -H "Content-Type: application/json" -d '{"test-var": "test-key"}'
Key-val correctly write in redis!
# curl -X GET localhost:8089/get_key?key=test-var
test-key
# curl -X DELETE localhost:8089/del_key -H "Content-Type: application/json" -d '{"test-var": "test-key"}'
Key correctly deleted!
# curl -X GET localhost:8089/get_key?key-test-var
404 page not found

```
