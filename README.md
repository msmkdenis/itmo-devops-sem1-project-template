# Финальный проект 1 семестра

REST API сервис для загрузки и выгрузки данных о ценах.

## Требования к системе

Go - 1.23.3

Postgresql - 15

Git - 2.34 и выше

curl - для тестирования api

Linux (Ubuntu 22.04 и выше)

Процессор: x86_64 (64-битный)

Оперативная память: минимум 2 ГБ

Свободное место на диске: минимум 100 МБ

## Установка и запуск

Скачать и установить postgresql: https://www.postgresql.org/download/

Склонировать проект:

```
git clone git@github.com:msmkdenis/itmo-devops-sem1-project.git
```

Перейти в директорию проекта и выполнить в терминале

создать бд и пользователя
```
sudo su postgres
psql
\i create_db.sql
exit
exit
```

собрать проект

```
go mod tidy
go build -o app cmd/main.go
```

запустить проект

```
./app
```

Сервер будет доступен по адресу: http://localhost:8080

## Тестирование

Для тестирования api можно воспользоваться curl и выполнить запросы:

Загрузить данные из архива:

```
curl -s -F "file=@sample_data.zip" "http://localhost:8080/api/v0/prices"
```

Выгрузить данные в виде архива data.zip:

```
curl -v -o data.zip "http://localhost:8080/api/v0/prices"
```

## Контакт

[Напишите в tg](https://t.me/kmsdenis)
