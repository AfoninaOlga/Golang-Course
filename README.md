# Golang-Course

## Демо

![demo](./images/demo.gif)

## Использование

### Сборка

`make` или `make all` собирает xkcd-server и web-server
Для сборки исполняемых файлов по отдельности можно запустить `make xkcdserver` и `make webserver` соответственно

### Тесты

`make test` запускает тесты и записывает результат покрытия в `coverage.html`

`make e2e` запускает end-to-end тест xkcd-сервера

`make lint` запускает `golang-ci-lint`

`make sec` запускает `trivy` и `govulncheck`

### Запуск

`./xkcd-server` или `./xkcd-server -c <path to config file> -p <port>`

> [!WARNING]
> Изменив порт xkcd-сервкра, не забудьте изменить адрес `api_url` в файле коннфигурации, передаваемом web-серверу

`./web-server` или `./web-server -c <path to config file> -p <port>`

## xkcd-server

Порт по умолчанию &mdash; 8080
Сервис, который по REST API выдает по поисковым словам url картинок, a также принимает команду на обновление базы.
- `POST /update` возвращает нам HTTP code ОК после обновления, количество new/total comics в виде JSON
- `GET /pics?search="I'll follow your questions"` возвращает урлы картинок в JSON формате
- `POST /login` с телом `{"name": "<name>", "password":"<password>"}` возвращает JWT токен
- `POST /register` с телом `{"name": "<name>", "password":"<password>"}` добавляет пользователя (не администратора) в базу

**Админ**: `{"name": "admin", "password": "admin"}`

```cmd
curl --request POST --data '{"name": "admin", "password": "admin"}' localhost:8080/login
```

**Добавление пользователя:**

```cmd
curl --request POST --data '{"name": "user1", "password": "pass123"}' localhost:8080/register
```

**Получение токена:**

```cmd
curl --request POST --data '{"name": "user1", "password": "pass123"}' localhost:8080/login
```

**Пример запроса с токеном:**

```cmd
curl -H 'Authorization: <token>' -v localhost:8080/pics?search="apple,doctor"
```

С сайта http://xkcd.com  скачивается описание всех комиксов, затем нормализуются слова. Хранятся В базе данных следующим образом:
- `Comics(id, url)`
- `Keywords(id, word)`
- `ComicsKeywords(id, comic_id, keyword_id)` &mdash; связь комиксов и ключевых слов
  И представление `Indexes(comic_id, keyword)` &mdash; `join` таблиц `Keywords` и `ComicsKeywords` 

## web-server

Порт по умолчанию &mdash; 8081
Веб-сервер, генерирующий html для xkcd-server
- `GET /login` генерирует страницу с формой для входа
- `POST /login` записывает токен в cookie и перенаправляет на `/comics`
- `GET /comics` генерирует страницу поиска
- `GET /comics?search=<keywords>` генерирует страницу с результатами поиска по <keywords>

## Конфигурация

Путь до файла конфигурации можно задать флагом `-c`, порт &mdash; флагом `-p`
В корневой директории лежит `config.yaml` файл, в котором пользователь может задать параметры:
  - `source_url` &mdash; источник для загрузки комиксов (по умолчанию https://xkcd.com)
  - `parallel` &mdash; число горутин для запроса комиксов
  - `update_time` &mdash; время ежедневного обновления базы комиксов
  - `api_port` &mdash; порт xkcd-сервера
  - `web_port` &mdash; порт web-сервера
  - `dsn` &mdash; строка подключения для базы комиксов
  - `concurrency_limit` &mdash; число запросов, обрабатываемых xkcd-сервером единовременно
  - `rate_limit` &mdash; максимальное число запросов от пользователя в секунду
  - `api_url` &mdash; адрес xkcd-сервера
  - `token_duration` &mdash; время жизни JWT токена в минутах
