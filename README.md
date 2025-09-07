```json
internal/
  application/        // use-cases + ports (интерфейсы)
    ports/            // контракты для хранилищ/сервисов
    usecases.go
  domain/             // сущности и доменные правила (файлы прямо тут)
    employee.go
    company.go
  infrastructure/     // реализации портов
    postgres/
    s3/
    cache/
  delivery/           // транспорт (HTTP/gRPC)
    http/
      handlers.go
      middleware/
  initialization/     // wiring + конфиг
  migrations/         // goose миграции
configs/              // конфиги
pkg/                  // технические утилиты
```

Архитектура (минимум, но правильно)

#### Слои и роли
**domain/** — сущности и правила предметки (value-objects, инварианты). Никаких импортов из фреймворков и БД.
**application/** — сценарии (use-cases) и ports/ (интерфейсы к внешнему миру: репозитории, транзактор, мейлер и т.п.).
**infrastructure/** — реализации портов (Postgres/S3/Redis и т.д.). Никакой бизнес-логики.
**delivery/** — транспорт (HTTP/gRPC/CLI): парсит вход, зовёт use-case, форматирует ответ.
**initialization/** — wiring: конфиг, DI, старт/грейсфул-шатдаун.
**migrations/, configs/, pkg/** — техничка; 
**pkg/** только для утилит без предметной области.

#### Направление зависимостей

```json
delivery  →  application(ports,usecases)  →  domain
↑
infrastructure (implements ports)
```

Domain ни от кого не зависит. Application знает только domain и свои интерфейсы. Infrastructure знает domain (для маппинга сущностей) и application/ports (чтобы реализовать контракты).
Use-case как контракт
Один UC = один тип с методом Handle(ctx, In) (Out, error) либо без In/Out, если тебе нужно совсем «голое» API.
Внутри UC — только orchestration: вызовы портов, доменные операции, транзакции (через порт транзактора).
Никаких HTTP/SQL-деталей.
Интерфейсы (ports) — в application/ports/

Пример: EmployeeRepository, Transactor, Mailer.
Почему здесь: интерфейс лежит там, где используется; Application диктует «что нужно», Infra подстраивается.
Транзакции
Порт Transactor.WithinTx(ctx, func(ctx) error) error.
UC, которому нужна запись, оборачивает логику в WithinTx. Реализация живёт в Infra (pgx/sqlx), а UC об этом не знает.

#### Ошибки

Доменные ошибки типизируй в domain (например, ErrEmployeeNotFound).
Infra маппит драйверные/SQL-ошибки → в доменные или технические (для логов). Delivery уже решает HTTP-код.
Delivery: «тонкие» хендлеры
Разрешено: парсить вход, дернуть UC, отдать ответ.
Запрещено: бизнес-правила, SQL, транзакции, прямые вызовы Infra.
pkg/ — карантин утилит
Логгеры, форматтеры, request-id, таймеры — да.
Всё, что знает про предметку — нет (в domain/application).

#### Инициализация

```json
Считываешь конфиг → поднимаешь соединения → создаёшь реализации портов → собираешь UC → вешаешь на роуты.
```

Закрытие ресурсов (DB, S3) — в defer/graceful shutdown.

#### Тесты

Unit: domain и application (моки портов).
Integration: infrastructure (реальная БД/MinIO), флаги -tags=integration и отдельный DSN.
Эволюция без боли
Нужны отдельные файлы по сущностям — раскладываешь domain/employee.go, domain/company.go.
Много UC — делишь usecases.go на usecase_employee.go, usecase_process.go и т.д.
Понадобятся DTO/валидаторы/мапперы — добавляешь позже (но это уже кастом).