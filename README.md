Проект выполнен в рамках курса [Route256 2021 OZON](https://rutracker.org/forum/viewtopic.php?t=6201055)

#### Репозиторий с заданиями
https://github.com/ozonmp/omp-docs

# О проекте
Репозиторий является частью проекта [Logistic-pack-api](https://github.com/StormBeaver/logistic-pack-api)

## logistic-pack-facade

Сервис читает сообщения в следующих кафка-топиках: 
  - "created";
  - "updated";
  - "removed".

и выводит сообщения в `stdout` в человекочитаемом виде.

### Топики "created", "updated", "removed"

Сообщения в топиках "created" сериализованы в json. 
Сообщения десериализуются в объект с помощью `json.Unmarshal`, затем, хранимый в jsonb payload повторно десереализуется при помощи `internal/model/parsePackEvent`, после чего сообщение выводятся в `stdout`.

## Docker

Описан докерфайл для образа фасада.

## Makefile

В мейкфайле описаны команды для локального запуска приложения, сборки докер-образа, и его запуска.