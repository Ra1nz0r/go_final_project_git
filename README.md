<div align="center"> <h1 align="center"> Планировщик - выпускной проект. </h1> </div>

- __В данном проекте выполнены все задания со звёздочкой__:
    - [x] Возможность определять порт, если есть переменная окружения TODO_PORT.
    - [x] Возможность определять путь к файлу базы данных, если есть переменная окружения TODO_DBFILE.
    - [x] Назначение задачи в указанные дни недели.
    - [x] Назначение задачи в указанные дни месяца, с возможностью указания конкретных месяцев.
    - [x] Реализован поиск[^1], независимо от регистра, по полям TITLE и COMMENT.
    - [x] Простейшая аутентификация, через переменную окружения TODO_PASSWORD.
    - [x] Реализован Dockerfile для создания образа.

__Программа выполняет функции планировщика, который возвращает следующую дату выполнения задачи в зависимости от переданных параметров её повторения.__

[Инструкция по локальному запуску и информация по приложению.](#local)\
[Инструкция по созданию и запуску Docker контейнера.](#docker)\
[Инструкция по запуску SQLite в терминале внутри запущенного контейнера.](#sqlite)

![logo](/internal/web/logo.gif)

<a name="local"></a>
***
#### Инструкция по локальному запуску и информация по приложению.

Пароль по-умолчанию для планировщика назначен: ```qwerty```\
По-умолчанию приложение запускается: ```0.0.0.0:7540```\
База создаётся c названием и путём: ```internal/storage_db/scheduler.db```\
Стандартные настройки для сервера лежат: ```internal/config/server_conf.go```\
Стандартные настройки для тестов лежат: ```internal/config/settings.go```\
Переменные окружения и их описание, хранятся в файле корня проекта: ```.env```\
Для генерации кода работы с базами данных, был использован: ```SQLC```[^2]

- Программу можно запускать двумя способами через терминал.
    - Обычные команды. 
    - Короткими командами из TaskFile.

_Для изменения стандартных параметров, как и пароля, нужно изменить ```TODO_PORT```, ```TODO_DBFILE``` или ```TODO_PASSWORD``` в ```.env``` файле корня проекта._
</div>

- ___Для запуска приложения в терминале.___\
```go run ./...``` или ```task run```
- ___Для запуска тестов в терминале.___\
```go test -v ./... -count=1``` или ```task test```

___!!! ЗАПУСК ТЕСТОВ УДАЛЯЕТ ВСЕ ДАННЫЕ ИЗ БАЗЫ ДАННЫХ !!!___

<a name="docker"></a>
***
#### Инструкция по запуску Docker контейнера.

- ___Для запуска сборки Docker.___\
```docker build -t scheduler_app:v1 .``` или ```task d_build```

- ___Обычный запуск с портом по-умолчанию, при изменении ```7540``` страница планировщика открываться не будет, сервер так и останется на ```7540``` порте.___\
```docker run --name="sched_app" -d -p 7540:7540 scheduler_app:v1```

- ___Для изменения стандартного порта необходимо передать ```TODO_PORT``` и указать его в ```-p```, тогда сервер будет доступен на переданном порте.___\
```docker run --name="sched_app" -e "TODO_PORT=7544" -d -p 7544:7544 scheduler_app:v1```

- ___Запустит контейнер на порт указанный в Dockerfile ```EXPOSE 7540```.___\
```docker run --name="sched_app" -d -P scheduler_app:v1```

- ___Для изменения стандартного пароля необходимо передать ```TODO_PASSWORD```, тогда стандартный пароль будет изменен и можно войти по своему.___\
```docker run --name="sched_app" -e "TODO_PASSWORD=yourPass" -d -p 7540:7540 scheduler_app:v1```

- ___Для изменения стандартного пути необходимо передать ```TODO_DBFILE```, тогда приложение создаст базу данных в соответствии с переданными, включая название базы.___\
```docker run --name="sched_app" -e "TODO_DBFILE=internal/yourPath/yourName.db" -d -p 7540:7540 scheduler_app:v1```

- ___Если необходимо изменить все стандартные значения. При изменении порта, обязательно передавать такое же значение через ```-p```.___\
```docker run --name="sched_app" -e "TODO_DBFILE=internal/yourPath/yourName.db" -e "TODO_PASSWORD=yourPass" -e "TODO_PORT=7544" -d -p 7544:7544 scheduler_app:v1```

- ___Запуск в интерактивном режиме, только при запущенном контейнере.___\
```docker exec -it sched_app /bin/bash```

<a name="sqlite"></a>
***
 #### Инструкция по запуску SQLite в терминале внутри запущенного контейнера

___Если имя таблицы и путь не изменялись, то запуск бинарника сразу откроет базу данных "scheduler.db". И можно работать с таблицей.___

___Для этого необходимо, при стандартных настройках:___
1. Запустить контейнер одним из обычных способов запуска [Docker](#docker) образа.
2. При запущенном контейнере выполнить ```docker exec -it sched_app /bin/bash``` команду.
3. Должна открыться консоль ```bash-5.1$ ```.
4. Выполнить команду ```./run_sqlite.sh``` и попадёте внутрь программы SQLite.
5. Для выхода написать ```.exit``` и снова ```exit```.

***

[^1]: Для поиска была сделана еще одна колонка в таблице, которая хранит данные из TITLE и REPEAT в нижнем регистре.
[^2]: После генерации нового кода SQLC необходимо заменить все типы значений ```ID``` ```int64``` на ```string``` в файле ```query.sql.go``` и ```models.go```. А также для методов ```CreateTask``` и ```UpdateTask```, изменить поля ```arg.Search``` на ```strings.ToLower(arg.Title+" "+arg.Comment)```.