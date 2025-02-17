# Тестовое задание для стажёра Backend-направления (зимняя волна 2025)

Использумые технологии:
* Язык Golang
* База данных PostgreSQL
* Swagger Server из Swagger Editor
* Локальный Swagger UI
* Docker

На данный момент реализована только бизнес логика API, все 4 запроса.
Были использованы почти все схемы (модели), кроме ErrorResponse.

## Инструкция к запуску docker контейнера.

В директории /docker лежат все необходимые файлы, нужно перейти в эту директорию и прописать:

docker-compose up

*По желанию можно переименовать общий контейнер (вместо docker на что-то другое), однако
необходимо сменить название директории в docker-compose.yaml:

`server:
build:
    context: ..
    dockerfile: <ИМЯ КОНТЕЙНЕРА>/Dockerfile`

Страница Swagger с реализованным наглядным API находится в:

http://localhost:8080/swagger/

## Локально:

В cmd/main.go напрямую лежат конфигурации для подключения к Postgres

dbClient, err := postgresql.NewClient(ctx, 5, dbHost, "5432", "postgres", "1234", "avito_shop")

Взять sql скрипт для базы данных можно по этим путям:

docker/init.sql
pkg/db_create_script.txt

Для запуска сервера локально:

cd "d:\Programming\Repositories\avito_winter2025\cmd\" ; if ($?) { go run main.go }

## Проблемы и сложности с проектом

1. Swagger

Об этой технологии я узнал полгода назад (летом 2024) во время выполнения задач практики, но не уложился и оставил его, выполнив самим базовую вёрстку с js скриптами.
Подробно спроектировав API и написав проект в одном yaml файле, можно автоматически сгенерировать сервер.
С подобным функционалом я встречался недавно - graphql.
После я запутался с путями, но разобрался и поставил Swagger UI.

2. PostgreSQL (схема)

Подключить базу данных было просто, потому что опыт с ним был. Но нужно было спроектировать таблицы для взаимодействия.
Создано три таблицы: users, transactions и inventories. Последняя таблица в начале должна была быть такой:
один пользователь - один инвентарь
Поэтому я попробовал хранить типы мерча и их количество в JSONB свойстве, но это было плохой идеей.
Из-за неудобства и траты большого количества времени я переделал эту таблицу. В итоге, один пользователь мог иметь до 10 сущностей в inventories (со своим типом мерча и его количеством).

3. Docker

С Docker был опыт полугодовой давности и представление о том, как нужно всё реализовать. Но этого не хватало, нужно было разбираться с подключением postgres (название не localhost, а db, отчего в main.go код дополнился условием) и swagger (смена путей). Однако, задача была выполнена и протестирована вручную.
