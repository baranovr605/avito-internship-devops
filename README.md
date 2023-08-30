# Avito тестовое задание для DevOps (2023)

## Содержание
- [Выполненные задания](#finish_tasks)
- [Структура и описание проекта](#project_descrption)
- [Описание конфигурации Golang приложения](#config_golang)
- [Описание конфигурации Nginx](#config_nginx)
- [Описание конфигурации Redis](#config_redis)
- [Описание работы с security.sh](#automaticly_security)
- [Настройка пароля вручную](#manual_security)
- [Запуск проекта в docker-compose](#run_project)

## <a name="finish_tasks"></a>Выполненные задания
- [x] В качестве базового образа используется Debian (используется debian:12.1-slim)
- [x] Развертывание приложения при помощи docker-compose версии 3.3 или старше (написан docker-compose и Dockerfile для всех сервисов)
- [x] Реализовано приложение на Golang
- [x] Работа в Redis происходит со строками (для удаления по ключу используется DEL)
- [x] На Redis поддержана аутентификация (аутентификация по заданному паролю и пользователю с помощью ACL)
- Redis и приложение общаются по зашифрованному каналу (TLS-соединение).


## <a name="project_descrption"></a>Структура и описание проекта
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
│   ├── users.acl
│   └── redis.conf
├── security.sh
└──README.md
```

- app/ - директория с Golang приложением 
- nginx/ - директория с конфигурацией и Dockerfile для запуска Nginx в контейнере 
- redis/ - директория с конфигурацией и Dockerfile для запуска REDIS в контейнере и файлом ACL 
- docker-compose - compose файл для запуска всех сервисов в контейнерах 
- security.sh - shell скрипт, используемый для генерации файла с паролем для Golang и конфига ACL для Redis 

В качестве базового образа используется debian:12.1-slim во всех сервисах \
Приложение запускается на порту 8089 \
Запросы проксируются с nginx на Golang приложение, которое делает непосредственно запросы к Redis 

Приложение использует 3 эндпоинта: 
1. POST /set_key { "[key]": "[val]" } - добавление ключ-значение в REDIS
2. GET /get_key?key=[key] - добавление ключ-значение в REDIS
3. DELETE /del_key { "key": "[key]" } - удаление ключ-значение по ключу

## <a name="config_golang"></a>Описание конфигурации Golang приложения
Приложение Golang конфигурируется с помощью environment в docker-compose.yml
Список переменных окружения:
1. APP_PORT - порт приложения Golang. Пример задания ":8000". В начале обязательно указывать двоеточие. При изменении порта так же рекомендуется изменить значение EXPOSE в Dockerfile.goapp с Golang
2. REDIS_HOST - хост Redis с портом, на котором работает Redis. Пример: "redis-go:6379". Указывается в формате "хост:порт"
3. REDIS_USER - пользователь Redis, который используется для выполнения запросов к Redis. Пользователь по умолчанию "default" отключен. Подробнее про конфигурирование в "Описание работы с security.sh"
4. REDIS_PASS_FILE - путь до файла с паролем к Redis внутри контейнера. По умолчанию путь к файлу, смонтированным через docker secret. Подробнее про конфигурирование в "Описание работы с security.sh"

## <a name="config_nginx"></a>Описание Nginx конфигурации
Nginx проксирует все запросы по 3 endpoint к приложению Golang \
Обращение идет по DNS имени (имени сервиса) и объявлено в директиве upstream в nginx/nginx.conf \
По умолчанию используется имя сервиса goapp и следующая конфигурация: 
```
    upstream backend {
        server goapp:8000;
    }
```
В случае, если необходимо запустить приложение локально, либо будет использоваться другой сервер с приложением, необходимо заменить goapp на другой адрес \
Для изменения порта работы Nginx внутри контейнера необходимо внести правки в следующую строку:
```
listen 80 default_server;
```
80 - значение по умолчанию. При изменении стоит изменить EXPOSE в Dockerfile.nginx

## <a name="config_redis"></a>Описание Redis конфигурации
Обращения от Golang приложения идут к Redis \
Основные параметры для изменения в redis/redis.conf
1. bind - адрес, к которому будет привязан Redis 
2. port - порт, на котором функционирует REDIS сервер (По умолчанию: 6379). При изменении рекомендуется изменить EXPOSE в Dockerfile.redis
3. aclfile - путь до ACL файла с именем и паролем для пользователя Redis внутри контейнера

## <a name="automaticly_security"></a>Описание работы с security.sh
Для корректной работы security.sh требуется установленный sha256. Запуск производить на OC Linux. \
security.sh предназначен для генерации файла с паролем пользователя для приложения Golang и для генерации ACL файла для Redis. \
Причины использования ACL для Redis вместо стандартного пользователя с паролем:
1. С помощью ACL в дальнейшем можно более детально настраивать доступы для пользователей Redis
2. С помощью ACL в файле можем хранить пароль в формате SHA256, что позволяет увеличить безопасность

После генерации файл с паролем для Golang и Redis монтируются в контейнеры с помощью secret, что позволяет избежать наличие пароля в args контейнера или в переменных окружения. \
Работа с security.sh: 
Пример генерации файла ACL и файла с паролем для приложения Golang
```
./security.sh gen_pass admin pass
```
- admin - имя пользователя, который будет иметь доступ в Redis
- pass - Пароль пользователя, который будет использован создаваемым пользователем в Redis
После работы завершения работы скрипта будет дополнен файл redis/users.acl пользователем с паролем в формате SHA256 и будет создан файл app/RedisPass с паролем для Redis

## <a name="manual_security"></a>Ручная настройка пароля
Если необходимо произвести ручную настройку пароля и доступа пользователя Redis необходимо сделать следующее:

1. Изменить файл users.acl добавив пользователя, его доступ к Redis и пароль (рекомендуется в формате SHA256)
2. Создать файл RedisPass в директории app и записать пароль к пользователю, указанному в redis/users.acl
3. Изменить пользователя для Golang приложения в переменных окружения в docker-compose

После внесения данных правок, можно запустить проект без помощи security.sh 
Монтирование паролей так же будет происходит через secret, что безопаснее, чем через переменные окружения и args

## <a name="run_project"></a>Запуск проекта в docker-compose
Для запуска проекта локально необходимо:

1. Установленный Docker, docker-compose 
Если данное ПО не установлено, необходимо обратиться к официальной документации: https://docs.docker.com/engine/install/

2. Установленный Git. Если не установлен, обратиться к официальной документации: https://git-scm.com/downloads

3. Установленный curl для проверки корректности работы приложения (опционально)


Далее будет приведен пример запуска приложения на Linux, используемый дистрибутив: Ubuntu 20.04 

Клонируем проект с Git и перейдем в директорию с проектом:
```
git clone https://github.com/baranovr605/avito-internship-devops.git
cd avito-internship-devops
```

Сконфигурируем пароли:
```
./security.sh gen_pass admin pass
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