# Тестовое задание на позицию стажера-бекендера
В этом репозитории находится реализация [тестового задания](https://github.com/avito-tech/autumn-2021-intern-assignment)
на позицию стажёра-бекендера в Avito.
#### Что выполнено?
- [x] Основное задание
- [x] Доп. задание №1
- [x] Доп. задание №2
- [x] Докеризация
- [x] Написаны тесты
- [x] Человеко-читабельные описания ошибок API
## Запуск
```
git clone https://github.com/alien-agent/users-balance-microservice

cd users-balance-microservice
```
Для запуска API сервера используется следующая команда:
```
docker-compose up
```
Для запуска тестов:
```
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```
## Детали реализации
Микросервис написан на Go, СУБД - PostgreSQL. Все запросы/ответы принимаются/отдаются в формате JSON.

#### Конфигурация
Некоторые параметры сервиса можно настраивать с помощью файлов конфигурации. Доступны следующие параметры:
 - `server_port` - порт, на котором API сервер будет принимать запросы
 - `rates_expiration` - время истечения "свежести" (частота обновления) курсов обмена валют
 - `dsn` - строка с настройками подключения к БД PostgreSQL

По умолчанию используется файл конфигурации `dev.yml`, а при запуске внутри Docker - `local.yml`. Также возможна 
конфигурация с помощью переменных среды - их приоритет выше, чем у файлов конфигурации. Соответствующие переменные среды
должны иметь префикс `APP_`, например: `APP_DSN`.


#### Как обеспечивается устойчивость данных о балансе пользователя?
1. На уровне *SQL Schema* установлен запрет на сохранение отрицательных балансов - даже в случае наличия бага в коде метод
получит ошибку базы данных и передаст ее внешнему сервису.
2. Все методы API, подразумевающие неоднократный доступ к базе данных (пример: денежный перевод - изменение баланса 
отправителя и получателя, создание транзакции в истории), вносят изменения в БД с помощью единственной транзакции. Пример:
во время денежного перевода деньги были списаны с отправителя, но при начислении получателю произошла ошибка - баланс отправителя
вернется в исходное состояние, API сообщит об ошибке.
3. Для представления баланса пользователя в API и БД используются 64-битные числа, что достаточно для представления любой
практической суммы денег.

#### Доп. задание №1
Конвертация валют происходит с использованием курса обмена валют с [бесплатного API](https://api.exchangerate.host/latest).
Если API обменных курсов недоступен, сервис API продолжает работать в штатном режиме, выдавая ошибку только при запросе 
баланса в отличной от рубля валюте. Курсы валют кешируются на 10 минут (время можно настроить в файле конфигурации - 
параметр `rates_expiration`). 

#### Доп. задание №2


## Project Layout

The starter kit uses the following project layout:

```
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── config               configuration files for different environments
├── internal             private application and library code
│   ├── config           configuration library
│   ├── deposit          deposit-related features
│   ├── entity           database models
│   ├── errors           error types and handling
│   ├── test             helpers for testing purpose
│   └── transaction      transaction-related features
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── log              structured and context-aware logger
└── testdata             test data scripts
```

Этот проект основан на [Go RESTful API Starter Kit](https://github.com/qiangxue/go-rest-api).
Огромная благодарность [@qiangxue](https://github.com/qiangxue) за его работу!