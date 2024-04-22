Запуск командой из корня проекта: 
```
go run cmd/main.go
```

Поддерживаются флаги, переменные окружения (переменные имеют приоритет)

- `-s` (SERVER_ADDRESS) сервер (default "localhost:8080")
- `-w` (FILE_WORKERS) количество обработчиков файлов (default 5)

Доступен `GET` запрос с параметрами: 

- `dir` - директория с файлами
- `word` - искомое слово

Протестировать можно запустив приложени и выполнив curl запрос

```
curl -G -v -d "dir=./examples" -d "word=World" http://localhost:8080/files/search
```

В ответ придет сообщение c заголовком `X-Cache: None`, т.к. к данной директории запрос выполнялся впервые.

```
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /files/search?dir=./examples&word=World HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Cache: None
< Date: Mon, 22 Apr 2024 07:05:56 GMT
< Content-Length: 18
< 
["file1","file3"]
* Connection #0 to host localhost left intact
```
При повторном запросе заголовок изменится на `X-Cache: Cached`

```
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /files/search?dir=./examples&word=World HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Cache: Cached
< Date: Mon, 22 Apr 2024 07:08:06 GMT
< Content-Length: 18
< 
["file1","file3"]
* Connection #0 to host localhost left intact

```

Задача

Необходимо реализовать компонент для поиска ключевого слова в файлах. Файлы содержат одно или несколько предложений.
Примеры файлов в каталоге examples.

Результат

Веб-сервис, который содержит один метод HTTP GET files/search. Метод возвращает список имен файлов в формате JSON, которые содержат это слово.

Требования

Допускается решение на одном из языков программирования: C#, Java, Python, JS, C/C++.
Для реализации поиска разрешается использовать только стандартные библиотеки, пакеты для вашего языка программирования.
Для реализации веб-сервиса можно подключать сторонние пакеты.
Тесты, комментарии и пояснения к коду приветствуются.

*Для решения на go
Версия 1.21+
Если слово не найдено, то возвращаем nil для списка файлов. 
Если при поиске произошла ошибка, то возвращаем ошибку и nil для списка файлов.
Также необходимо добавить тесты для следующих случаев: слово не найдено; при обработке файла возникла ошибка.

Плюсами будут.
Реализация на go. Заготовка для go в каталоге pkg
Поиск слова за O(1).
Реализация параллельного поиска по файлам.

