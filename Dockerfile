ARG PORT=8080

FROM golang:1.25-alpine

WORKDIR /app

# Модули
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY ./cmd ./cmd/
COPY ./internal ./internal/


# Сборка
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/app/sem1-final-project-hard-level/main.go

# Создание пользователя
RUN adduser -D -s /bin/sh -u 2000 appuser  && \
    chown -R appuser:appuser /app

# Переключение на пользователя
USER appuser

ARG PORT
ENV PORT=${PORT}
EXPOSE ${PORT}

CMD [ "./main" ]

