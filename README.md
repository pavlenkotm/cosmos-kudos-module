# Cosmos Kudos Module

Модуль для Cosmos SDK-цепочки, добавляющий функционал отправки "кудосов" (очков благодарности) между адресами.

## Описание

Cosmos Kudos Module — это минимальный пример модуля для Cosmos SDK, демонстрирующий создание простого социального токена без сложной экономики. Модуль позволяет пользователям отправлять друг другу очки благодарности (kudos) с опциональными комментариями.

## Возможности

- Отправка кудосов между адресами
- Проверка баланса кудосов для любого адреса
- Просмотр таблицы лидеров (топ получателей кудосов)
- История всех транзакций кудосов
- CLI команды для взаимодействия с модулем
- gRPC/REST API для запросов

## Технологический стек

- **Язык**: Go 1.21+
- **Фреймворк**: Cosmos SDK v0.50+
- **Protobuf**: Для определения типов сообщений и запросов

## Структура проекта

```
cosmos-kudos-module/
├── x/kudos/                    # Основной модуль
│   ├── keeper/                 # Keeper для управления состоянием
│   │   ├── keeper.go          # Основная логика keeper
│   │   ├── msg_server.go      # Обработчик сообщений
│   │   ├── query_server.go    # Обработчик запросов
│   │   └── keeper_test.go     # Тесты keeper
│   ├── types/                  # Типы данных модуля
│   │   ├── keys.go            # Ключи для KVStore
│   │   ├── errors.go          # Ошибки модуля
│   │   ├── codec.go           # Регистрация кодеков
│   │   ├── msgs.go            # Сообщения и валидация
│   │   ├── genesis.go         # Genesis состояние
│   │   ├── tx.pb.go           # Сгенерированные типы транзакций
│   │   └── query.pb.go        # Сгенерированные типы запросов
│   ├── client/cli/             # CLI команды
│   │   ├── tx.go              # Команды транзакций
│   │   └── query.go           # Команды запросов
│   └── module.go               # Регистрация модуля
├── app/                        # Пример интеграции
│   └── app.go                 # Пример приложения с модулем
├── proto/kudos/                # Protobuf схемы
│   ├── query.proto            # Схема запросов
│   └── tx.proto               # Схема транзакций
├── go.mod
└── README.md
```

## Модель данных

### KudosBalance

Хранит количество полученных кудосов для каждого адреса:

- **Ключ**: `KudosBalancePrefix + address`
- **Значение**: `uint64` (количество кудосов)

### KudosHistory

Хранит историю всех транзакций кудосов:

- **Поля**:
  - `from_address` — адрес отправителя
  - `to_address` — адрес получателя
  - `amount` — количество кудосов
  - `comment` — комментарий (до 140 символов)
  - `timestamp` — временная метка

## Сообщения

### MsgSendKudos

Отправка кудосов от одного адреса другому.

**Поля**:
- `from_address` (string) — адрес отправителя
- `to_address` (string) — адрес получателя
- `amount` (uint64) — количество кудосов
- `comment` (string) — комментарий (максимум 140 символов)

**Правила валидации**:
- Отправитель не может отправить кудосы самому себе (`from_address` != `to_address`)
- Количество должно быть больше 0 (`amount` > 0)
- Длина комментария не должна превышать 140 символов

## gRPC/REST API

### Запросы

#### QueryKudosBalance

Получить баланс кудосов для адреса.

**Запрос**:
```protobuf
message QueryKudosBalanceRequest {
  string address = 1;
}
```

**Ответ**:
```protobuf
message QueryKudosBalanceResponse {
  uint64 balance = 1;
}
```

**REST**: `GET /kudos/balance/{address}`

#### QueryKudosLeaderboard

Получить топ N получателей кудосов.

**Запрос**:
```protobuf
message QueryKudosLeaderboardRequest {
  uint32 limit = 1;
}
```

**Ответ**:
```protobuf
message QueryKudosLeaderboardResponse {
  repeated LeaderboardEntry entries = 1;
}
```

**REST**: `GET /kudos/leaderboard?limit=10`

## CLI команды

### Транзакции

#### Отправить кудосы

```bash
<appd> tx kudos send [to_address] [amount] --comment "Комментарий" --from [from_key]
```

**Пример**:
```bash
appd tx kudos send cosmos1abc... 10 --comment "Спасибо за ревью кода!" --from alice
```

### Запросы

#### Проверить баланс

```bash
<appd> query kudos balance [address]
```

**Пример**:
```bash
appd query kudos balance cosmos1abc...
```

#### Посмотреть таблицу лидеров

```bash
<appd> query kudos leaderboard [limit]
```

**Пример**:
```bash
appd query kudos leaderboard 10
```

## Интеграция в приложение

### Шаг 1: Добавить зависимость

```bash
go get github.com/pavlenkotm/cosmos-kudos-module
```

### Шаг 2: Импортировать модуль в `app.go`

```go
import (
    kudosmodule "github.com/pavlenkotm/cosmos-kudos-module/x/kudos"
    kudoskeeper "github.com/pavlenkotm/cosmos-kudos-module/x/kudos/keeper"
    kudostypes "github.com/pavlenkotm/cosmos-kudos-module/x/kudos/types"
)
```

### Шаг 3: Добавить store key

```go
keys := storetypes.NewKVStoreKeys(
    // ... другие ключи
    kudostypes.StoreKey,
)
```

### Шаг 4: Инициализировать keeper

```go
app.KudosKeeper = kudoskeeper.NewKeeper(
    appCodec,
    runtime.NewKVStoreService(keys[kudostypes.StoreKey]),
    logger,
)
```

### Шаг 5: Зарегистрировать модуль

```go
app.mm = module.NewManager(
    // ... другие модули
    kudosmodule.NewAppModule(appCodec, app.KudosKeeper),
)
```

### Шаг 6: Настроить порядок выполнения

```go
app.mm.SetOrderBeginBlockers(
    // ... другие модули
    kudostypes.ModuleName,
)

app.mm.SetOrderEndBlockers(
    // ... другие модули
    kudostypes.ModuleName,
)

app.mm.SetOrderInitGenesis(
    // ... другие модули
    kudostypes.ModuleName,
)
```

### Шаг 7: Зарегистрировать сервисы

```go
app.mm.RegisterServices(module.NewConfigurator(
    appCodec,
    app.MsgServiceRouter(),
    app.GRPCQueryRouter(),
))
```

## Тестирование

### Запуск тестов

```bash
go test ./x/kudos/...
```

### Тесты включают

- **Unit-тесты keeper**:
  - Отправка кудосов
  - Валидация (нельзя отправить самому себе)
  - Проверка баланса
  - Таблица лидеров
  - История транзакций

- **Тесты валидации сообщений**:
  - Проверка корректных адресов
  - Проверка количества
  - Проверка длины комментария

### Пример теста

```go
func TestSendKudos(t *testing.T) {
    k, ctx := setupKeeper(t)

    fromAddr := "cosmos1from"
    toAddr := "cosmos1to"

    err := k.SendKudos(ctx, fromAddr, toAddr, 100, "Great work!")
    require.NoError(t, err)

    balance := k.GetKudosBalance(ctx, toAddr)
    require.Equal(t, uint64(100), balance)
}
```

## Примеры использования

### Пример 1: Отправка кудосов за код-ревью

```bash
appd tx kudos send cosmos1reviewer... 5 \
  --comment "Спасибо за детальное ревью PR #123!" \
  --from developer \
  --chain-id kudos-chain-1
```

### Пример 2: Проверка баланса

```bash
appd query kudos balance cosmos1reviewer...
```

**Ответ**:
```json
{
  "balance": "25"
}
```

### Пример 3: Просмотр топ-10 лидеров

```bash
appd query kudos leaderboard 10
```

**Ответ**:
```json
{
  "entries": [
    {
      "address": "cosmos1top1...",
      "balance": "150"
    },
    {
      "address": "cosmos1top2...",
      "balance": "120"
    }
  ]
}
```

## Архитектурные решения

### Почему KVStore?

Модуль использует простое key-value хранилище для максимальной производительности и простоты. Балансы хранятся по ключу адреса, что обеспечивает O(1) доступ к данным.

### Почему отсутствует списание кудосов?

Модуль реализует одностороннюю систему благодарностей — кудосы могут только накапливаться. Это упрощает логику и предотвращает возможные споры о списании.

### История транзакций

История транзакций опциональна и хранится для аудита и аналитики. Каждая транзакция получает уникальный ID на основе глобального счетчика.

## Расширение модуля

### Возможные улучшения

1. **Периодическое обнуление балансов**: Добавить BeginBlocker для сброса балансов раз в эпоху
2. **Ограничение отправки**: Добавить лимит на количество кудосов, которые можно отправить за период
3. **Репутационная система**: Использовать кудосы для расчета репутации участников
4. **NFT награды**: Выдавать NFT за достижение определенных порогов кудосов
5. **Запросы истории**: Добавить gRPC запросы для получения истории транзакций

### Пример расширения: Добавление лимитов

```go
// В keeper.go
const MaxKudosPerDay = 100

func (k Keeper) SendKudos(ctx sdk.Context, from, to string, amount uint64, comment string) error {
    // Проверить дневной лимит
    dailySent := k.GetDailySent(ctx, from)
    if dailySent + amount > MaxKudosPerDay {
        return ErrDailyLimitExceeded
    }

    // ... существующая логика
}
```

## Производительность

- **Отправка кудосов**: O(1) — простое обновление значения в KVStore
- **Проверка баланса**: O(1) — прямой доступ по ключу
- **Таблица лидеров**: O(n log n) — требует итерации и сортировки всех балансов

## Безопасность

### Валидация входных данных

Все сообщения проходят валидацию в `ValidateBasic()` перед выполнением:
- Проверка формата адресов
- Проверка корректности значений
- Проверка длины строк

### Защита от спама

Рекомендуется добавить:
- Минимальную комиссию за транзакции
- Rate limiting для предотвращения спама
- Ограничения на количество транзакций от одного адреса

## Лицензия

MIT

## Контакты

Для вопросов и предложений создавайте issue в GitHub репозитории.

## Дополнительные ресурсы

- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network/)
- [Building Modules Guide](https://docs.cosmos.network/main/building-modules/intro)
