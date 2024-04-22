Запуск командой из корня проекта: 
- `go run cmd/main.go`

Поддерживаются флаги, переменные окружения (переменные имеют приоритет)

- `-s` (SERVER_ADDRESS) сервер (default "localhost:8080")
- `-w` (FILE_WORKERS) количество обработчиков файлов (default 5)

Доступен `GET` запрос с параметрами: 

- `dir` - директория с файлами
- `word` - искомое слово

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

