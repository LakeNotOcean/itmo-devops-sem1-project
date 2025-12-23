# Финальный проект 1 семестра

REST API сервис для загрузки и выгрузки данных о ценах.

## Требования к системе

Предварительно необходимо задать переменные окружения в файле `configs/.env` в соответствии с шаблоном `configs/env.example`:
- Переменные разделов *Backend* и *Database* обязательны;
- Переменные разделов *Docker* и *Deploy* необходимы для запуска на удаленном сервере;

### Для локального запуска

- Go 1.25 или выше;
- Postgres 18 или выше;

### Для запуска с Docker

- Docker версии 28 или выше;
- Docker Compose версии v2 или выше;

### Для развертывания в Yandex Cloud:

- Установленная утилита `yc` актуальной версии;
- Пара ssh-ключей для доступа к виртуальной машине;
- Аккаунт с доступом к Yandex Cloud;
- Утилита `jq` актуальной версии;

## Установка и запуск

### Обязательные шаги

1. Склонировать репозиторий: `https://github.com/LakeNotOcean/itmo-devops-sem1-project.git`.
2. Настроить переменные окружения `configs/.env`.

### Установка и запуск локально

3. Создать пользователя и базу данных Postgres в соответствии с переменными окружения.
4. Выполнить установку зависимостей Go:
```
go mod vendor
```
5. Запустить приложение:
```
go run ./cmd/app/sem1-final-project-hard-level/main.go
```
6. Api будет доступно по адресу: `http://localhost:${PORT}/api/v0`

### Локальная установка с Docker

3. Установить актуальную версию Docker.
4. Выполнить сборку backend-сервиса c помощью скрипта:
```
scripts/prepare.sh
```
5. Запустить контейнеры:
```
sudo docker compose up --env-file=./configs/.env --build -d
```
6. Api будет доступно по адресу: `http://localhost:${PORT}/api/v0`

**Внимание!** Конфигурация Docker Compose включает установку PostgreSQL.

### Установка на удаленном сервере (Yandex Cloud)

Выполнить скрипт:
```
scripts/run.sh
```

При необходимости можно выполнить шаги отдельно:
- `scripts/helpers/create-infrastructure.sh` - настроить инфраструктуру на удаленном сервере.
- `scripts/helpers/install-dependencies.sh` - установить зависимости (Docker).
- `scripts/helpers/deploy-application.sh` - развернуть приложение на удаленном сервере.

**Внимание!** Если инфраструктура с указанными параметрами уже развёрнута, её перенастройка не производится.

## Тестирование

    ./scripts/tests.sh 1
    ![alt text](https://github.com/LakeNotOcean/itmo-devops-sem1-project/blob/main/docs/images/test-1.png)
    ./scripts/tests.sh 2
    ![alt text](https://github.com/LakeNotOcean/itmo-devops-sem1-project/blob/main/docs/images/test-2.png)
    ./scripts/tests.sh 3
    ![alt text](https://github.com/LakeNotOcean/itmo-devops-sem1-project/blob/main/docs/images/test-3.png)
    ./scripts/run.sh
    ![alt text](https://github.com/LakeNotOcean/itmo-devops-sem1-project/blob/main/docs/images/run-sh.png)
    GitHub Actions
    ![alt text](https://github.com/LakeNotOcean/itmo-devops-sem1-project/blob/main/docs/images/ci-cd.png)
    

## Контакты

- Telegram: @cheshirskins
- GitHub: https://github.com/LakeNotOcean