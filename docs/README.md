### Краткий обзор концепта с DSL настроек маршрутизации и `SEGMENTS` вместо жёсткой структуры сообщения

#### 1. Настройки маршрутизации

Ключевое отличие от предыдущей версии:

---

- вместо нескольких справочников таблица `address_book`/`address_dict`/`dict`/...(подобрать название) \
  Поля таблицы:
    - `id` int - идентификатор адреса (uuid/string/bigint/... на выбор)
    - `address` string - SITA адрес
    - `client_id` string - IATA код клиента или другой идентификатор
    - `status` string - активен/не активен, если вкратце
    - `type` string - [INTERNAL|EXTERNAL|DESTINATION|SERVICE|BAGGAGE_SPEC|...] (больше никаких адресов из переменных
      окружения)
    - `is_defualt` bool - для одного клиента(`client_id`) в рамках одного типа(`type`) может быть установленным как
      адрес по умолчанию
    - `delegated_to` string - `client_id`, на кого адрес делегирован

  Пример:
    ```json
      {
        "id": 1,
        "address": "XXXYYZZ",
        "client_id": "AB",
        "status": "ACTIVE",
        "type": "INTERNAL",
        "is_default": false,
        "delegated_to": "CD"
      }
    ```

---

- для справочника добавляем таблицу `properties`/`address_properties`, где 1 строка является 1 значением \
  Поля:
    - `address_id` int - ссылка на адрес
    - `name` string - название свойства
    - `value` string - значение свойства

  Пример:
  ```json
    [
      {
        "address_id": 1,
        "name": "email",
        "value": "example1@email.com"
      },
      {
        "address_id": 1,
        "name": "email",
        "value": "example2@email.com"
      },
      {
        "address_id": 1,
        "name": "protocol",
        "value": "SMTP"
      }
    ]
  ```
  Таким образом, вместо поля `jsonb` словарик свойств

---

- `ПРАВИЛА`. Вместо таблиц нескольких таблиц правил 1 таблица `rules` со следующими полями:

    - `id` int - идентификатор правила
    - `client_id` string - идентификатор клиента, которому оно принадлежит
    - `context` string - определение, для чего используются (например, для отправки багажных без заголовка, приём с
      внешки и т.п.)
    - `conditions` jsonb - `JSON` с условиями выполнения правила
    - `links` []string/[]int - ссылка на строки в таблице `address_dict`

В зависимости от одного из 3‑х вариантов(а их может быть и больше вплоть до написания небольшого ЯП) примеры будут
разниться.

1) Когда все условия выполняются как `AND`, тогда `conditions` будут проще:

    ```json
    [
      {"source": "header", "type": "receiver", "op": "eq", "value": "XXXYYZZ"},
      {"source": "body", "type": "departure_airport", "op": "in", "value": ["SU", "SVO"]}
    ]
    ```
   В таком случае мы получаем условие: если `@header.receiver` == `XXXYYZZ` И `@body.departure`.Contains(`SU`, `SVO`),
   правило будет считаться выполненным.

2) Когда мы определяем сверху над всеми элементами тип соединения (AND || OR):

    ```json
    {
      "union": "OR",
      "conditions":     [
        {"source": "header", "type": "receiver", "op": "eq", "value": "XXXYYZZ"},
        {"source": "body", "type": "departure_airport", "op": "in", "value": ["SU", "SVO"]}
      ]
    }
    ```

   В таком случае мы получаем условие: если `@header.receiver` == `XXXYYZZ` ИЛИ `@body.departure`.Contains(`SU`, `SVO`),
   правило будет считаться выполненным.

3) Когда мы группируем в блоки наши условия:
    ```json
    {
      "statements_union": "OR",
      "statements":     [
        {
          "union": "OR",
          "conditions":     [
            {"source": "header", "type": "receiver", "op": "eq", "value": "XXXYYZZ"},
            {"source": "body", "type": "departure_airport", "op": "in", "value": ["SU", "SVO"]}
          ]
        },
        {
          "union": "AND",
          "conditions":     [
            {"source": "header", "type": "sender", "op": "eq", "value": "GGGOOPP"},
            {"source": "body", "type": "flight_number", "op": "neq", "value": "SS77413"}
          ]
        } 
      ]
    }
    ```
   В таком случае мы получаем условие: если (`@header.receiver` == `XXXYYZZ` ИЛИ `@body.departure`.Contains(`SU`, 'SVO')
   И (`@header.sender` == `GGGOOPP` И `@body.flight_number` != `SS77413`)
   правило будет считаться выполненным.

---

#### 2. Массив сегментов телеграммы вместо жёсткой Go-структуры.

На данный момент в сервисе `Router` мы принимаем жёсткую структуру `go` с фиксированным набором полей, однако, \
если мы определяем такие блоки телеграммы как отправитель, получатель, аэропорт вылета, тип судна и т.д. как сегменты, \
то можем, во-первых, свести эти свойства в структур `Segment`:

```go
package segment

type Segment struct {
	Source string
	Type   string
	Value  string
}
```

Либо же интерфейс

```go
package segment

type Segment interface {
	Source() string
	Type() string
	Value() string
}
```

В таком случае, если ожидать на входе в роутер телеграмму с минимальным набором полей метаданных (`client_id`,
`from_sita`), \
а остальные данные свести к массиву сегментов, то, во-первых, сегменты спокойно "ложатся" на структуру `conditions` по
`source` и `type`, \
во-вторых, позволяет нам лишь расширить структуру передаваемых сегментов только в 1ом месте, если парсер/language server
будут нам СРАЗУ \
возвращать сегменты, т.е., если изначально парсинг был не полным, а затем мы решили его расширить, \
тогда нам нужно лишь расширить парсер и поддержать новое поле для настройки(condition) во Frontend, \
чтобы пользователи могли пользоваться новой настройкой (однако, мы уже можем руками в бд занести новые условия, если это
нужно было "ещё вчера").