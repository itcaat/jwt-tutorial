# jwt-tutorial

```
jwt-tutorial/
jwt-tutorial/
│── ui/
│   ├── index.html        # Форма логина (Bootstrap)
│   ├── dashboard.html    # Личный кабинет (Bootstrap)
│   ├── app.js           # JS-код (jQuery + AJAX)
│   └── README.md        # Документация по UI
│── jwt-app/                  # Go-приложение (Backend API)
│   │── cmd/                  # Главная точка входа (main.go)
│   │   └── main.go
│   │
│   ├── internal/             # Внутренний код (логика API)
│   │   ├── handlers/         # HTTP-обработчики
│   │   │   ├── auth.go       # Авторизация и генерация JWT
│   │   │   ├── protected.go  # Защищённые маршруты (read/write)
│   │   │   └── home.go       # Главная публичная страница
│   │   │
│   │   ├── middleware/       # Middleware-функции
│   │   │   ├── cors.go       # Проверка CORS
│   │   │   ├── jwt.go        # Проверка JWT
│   │   │   └── roles.go      # Проверка ролей
│   │   │
│   │   ├── models/           # Определение структур данных
│   │   │   ├── user.go       # Структура пользователя и claims
│   │   │   ├── request.go    # Запросы (логин и т.д.)
│   │   │   └── response.go   # JSON-ответы
│   │   │
│   │   ├── config/           # Конфигурация (переменные окружения)
│   │   │   └── config.go
│   │   │
│   │   ├── auth/             # Аутентификация
│   │   │   ├── jwt.go        # Работа с JWT (создание, проверка)
│   │   │   └── users.go      # Работа с пользователями (mock или БД)
│   │   │
│   │   └── routes/           # Маршрутизация API
│   │       └── router.go
│   │
│   ├── tests/                # Тесты (unit и integration)
│   │   ├── auth_test.go
│   │   ├── middleware_test.go
│   │   ├── protected_test.go
│   │   └── router_test.go
│   │
│   ├── go.mod                # Go модули (Go dependencies)
│   ├── go.sum                # Контрольная сумма зависимостей
│   ├── Dockerfile            # Docker-контейнер для API
│   ├── Makefile              # Сборка и запуск проекта
│   └── README.md             # Документация backend-а
│
├── ui/                       # Вынесенный UI (Frontend)
│   ├── index.html            # Главная страница (логин)
│   ├── dashboard.html        # Страница с доступом к API
│   ├── app.js                # Взаимодействие с API (авторизация, запросы)
│   ├── styles.css            # CSS-оформление
│   └── README.md             # Документация по UI
│
├── README.md                 # Основная документация проекта
```

# Описание логики

## cmd/main.go

- Создаёт HTTP-сервер на 8080.
- Использует mux.NewRouter() для управления маршрутами.
- Подключает маршруты через routes.RegisterRoutes(r).
- Выводит в консоль ссылку на сервер.
- Завершает работу с ошибкой, если сервер не запустится.

## internal/routes

### router.go
Это файл отвечает за регистрацию всех маршрутов API.

#### Публичные маршруты:
- / → Открытая главная страница.
- /login → Аутентификация (выдача JWT).

#### Защищённые маршруты (доступны только с токеном):
- /protected/read → Доступно для read и write.
- /protected/write → Доступно только для write.

## internal/handlers

### internal/handlers/home.go – Публичная главная страница
Этот обработчик отвечает за доступ к /.

- Отдаёт JSON с приветственным сообщением.
- Не требует аутентификации.

### internal/handlers/auth.go – Логин и выдача JWT
Этот обработчик отвечает за /login.

- Принимает username и password.
- Проверяет пользователя по mock-данным (можно заменить на БД).
- Создаёт JWT-токен с ролью (read или write).
- Отправляет токен в ответе.

### internal/handlers/protected.go – Защищённые маршруты
- /protected/read доступен всем авторизованным (read и write).
- /protected/write доступен только пользователям с ролью write.

## internal/middleware - проверка JWT и ролей.

### internal/middleware/cors.go – Проверка CORS
- Задает разрешающие заголовки для CORS

### internal/middleware/jwt.go – Проверка JWT
- Проверяет, что есть заголовок Authorization: Bearer <TOKEN>.
- Декодирует и валидирует JWT.
- Добавляет X-Username и X-Role в заголовки запроса для последующего использования.

### internal/middleware/roles.go – Проверка ролей
- Проверяет, что у пользователя есть необходимая роль.
- Если у пользователя read, но он пытается записать → 403 Forbidden.

## internal/models – Структуры

### internal/models/user.go – Структуры пользователей и JWT Claims

Этот файл определяет структуры для работы с пользователями и JWT.

- Claims – структура данных, хранящаяся в токене.
- Использует jwt.RegisteredClaims для стандартных полей (например, ExpiresAt).

### internal/models/request.go – Структуры входных данных
Этот файл определяет структуры для запросов, которые отправляет клиент.

- Credentials – используется в LoginHandler для обработки JSON-запроса.

### internal/models/response.go – Структуры ответов API
Этот файл определяет стандартные ответы API.

- AuthResponse – отправляется при успешном логине (возвращает token).
- ErrorResponse – используется для обработки ошибок (например, Unauthorized).


## Тестирование API

### internal/tests

```go test ./internal/tests```

#### Тестируем аутентификацию (auth_test.go)

- Логин с reader/password должен вернуть 200 OK
- В ответе должен быть token

#### Тестируем JWT Middleware (middleware_test.go)

- JWTMiddleware корректно валидирует токены.
- Запрос без токена должен вернуть 401 Unauthorized

#### Тестируем защищённые маршруты (protected_test.go)

- Тестируем доступ к /protected/read и /protected/write.
- reader может читать (200 OK)
- reader не может писать (403 Forbidden)

### Ручное тестирование API

Для ручного тестирования необходимо запустить приложение.

```bash 
$ cd jwt-app
$ go run cmd/main.go
# Ожидаемый ответ: Server running on http://localhost:8080
```

#### Тест 1: Проверка публичного маршрута (GET /)
```bash 
$ curl -X GET http://localhost:8080/
# Ожидаемый ответ: {"message": "Welcome to the public home page!"}
```

#### Тест 2: Логин как reader (POST /login)
```bash 
# Скопируй полученный токен, он нам понадобится для следующих запросов.
$ curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"reader","password":"password"}'
# Ожидаемый ответ: {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
```

#### Тест 3: Доступ к /protected/read с токеном
```bash 
# Заменяем <TOKEN> на полученный в предыдущем шаге.
$ curl -X GET http://localhost:8080/protected/read -H "Authorization: Bearer <TOKEN>"
# Ожидаемый ответ (Ошибка 403 Forbidden): {"message": "Protected data accessed!", "username": "reader", "role": "read"}
```

#### Тест 4: Попытка записи (POST /protected/write) с reader (ОШИБКА)
```bash 
$ curl -X POST http://localhost:8080/protected/write -H "Authorization: Bearer <TOKEN>"
# Ожидаемый ответ (Ошибка 403 Forbidden): {"error": "Forbidden: requires write access"}
```

#### Тест 5: Логин как writer
```bash 
$ curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"writer","password":"password"}'
# Ожидаемый ответ: {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
# Скопируй этот новый токен.
```

#### Тест 6: Доступ к /protected/write с writer (успешно)
```bash 
$ curl -X POST http://localhost:8080/protected/write -H "Authorization: Bearer <TOKEN>"
# Ожидаемый ответ: {"message": "Data successfully written!", "username": "writer"}
```
