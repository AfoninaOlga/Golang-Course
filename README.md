# Golang-Course

## Описание
Веб-сервис, который по REST API выдает по поисковым словам url картинок, a также принимает команду на обновление базы.
- `POST /update` возвращает нам HTTP code ОК после обновления, количество new/total comics в виде JSON
- `GET /pics?search="I'll follow your questions"` возвращает урлы картинок в JSON формате
- ``

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
- `ComicsKeywords(id, comic_id, keyword_id)` — связь комиксов и ключевых слов
  И представление `Indexes(comic_id, keyword)` — `join` таблиц `Keywords` и `ComicsKeywords` 

В корневой директории лежит `config.yaml` файл, в котором пользователь может задать параметры:
  - `source_url` — по умолчанию https://xkcd.com
  - `db_file` — по умолчанию `database.json`
  - `port` — порт сервера
- `make server` компилирует приложение сервера `xkcd-server`
- Порт задаётся флагом `-p` или опцией `port` в `config.yaml`
