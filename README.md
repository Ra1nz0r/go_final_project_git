- ___С чтением всех значений из файла енв.___ 
    ```docker run --name="sched_app" -d --env-file .env -p 7540:7540 scheduler_app:v1```

- ___Запустит на порт указанный в dockerfile EXPOSE 7540.___ 
    ```docker run --name="sched_app" -d -P scheduler_app:v1```

- ___Обычный запуск с портом по-умолчанию, при изменении не будет работать на введёном порту, сервер так и останется на 7540.___
```docker run --name="sched_app" -d -p 7540:7540 scheduler_app:v1```

- ___Если необходимо изменить стандартный порт.___
 ```docker run --name="sched_app" -e "TODO_PORT=7544" -d -p 7544:7544 scheduler_app:v1```

- ___Если необходимо изменить стандартный пароль___
```docker run --name="sched_app" -e "TODO_PASSWORD=gdfsd" -d -p 7540:7540 scheduler_app:v1```

- ___Если необходимо изменить стандартный путь к базе данных и название базы.___
```docker run --name="sched_app" -e "TODO_DBFILE=internal/sge_db/sdfgduler.db" -d -p 7540:7540 scheduler_app:v1```

- ___Если необходимо изменить все стандартные значения. При изменении порта, обязательно передавать такое же значение через -p.___
```docker run --name="sched_app" -e "TODO_DBFILE=internal/sge_db/sdfgduler.db" -e "TODO_PASSWORD=gdfsd" -e "TODO_PORT=7544" -d -p 7544:7544 scheduler_app:v1```

- ___Запуск в интерактивном режиме, только при запущенном контейнере:___
```docker exec -it sched_app /bin/bash```
