# 🛍️ Product Search Microservice

Микросервис **product-search** обеспечивает индексацию и полнотекстовый поиск товаров с поддержкой:

* **Массовой** и **инкрементальной** индексации через Kafka-сообщения
* **Фильтрации** по категориям, бренду и ценовым диапазонам
* **Сортировки** по релевантности, цене и популярности
* **Пагинации** и **подсветки** (highlighting) найденных фрагментов
* **Фасетной навигации** (aggregation buckets)
* **Metrics & Monitoring**: проверка статуса кластера и основных метрик

Построен по принципам **Clean Architecture**:

```
transport  → use_cases  → domain  → infrastructure
    ^                                    |
    └────────────── Docker & Kafka ──────┘
```

---

## 🚀 Возможности

* 🔄 **Bulk / Incremental Indexing** — через Kafka consumer
* 🔍 **Search API** с JSON-запросом (multi\_match, фильтры, sort, from/size)
* ✨ **Highlighting** ключевых слов в полях `name` и `description`
* 📊 **Facet Aggregations** по `category` и `brand`
* 📈 **Cluster Health** (`GET /product/health`)
* 📦 **Docker Compose** для Elasticsearch (кластер из 3-х нод) и Kibana
* 🧩 **Clean Architecture** (domain, use\_cases, transport, infrastructure)

---

## ⚙️ Технологии

* **Go**
* **gin-gonic/gin** (REST API)
* **github.com/elastic/go-elasticsearch/v8**
* **github.com/segmentio/kafka-go** (Kafka consumer)
* **Docker & Docker Compose**
* **Clean Architecture**

---

## 🛠 Makefile

В корне проекта есть `Makefile` с основными целями:

* `make help`
  Показать список целей.

* `make deps`
  Установить Go-модули и инструменты (swag, protoc-gen-go).

* `make build`
  Собрать бинарник в `bin/product-search`.

* `make run`
  Запустить сервис локально (сначала `build`, потом бинарник).

* `make docker-up` / `make docker-down`
  Поднять/остановить ES-кластер и Kibana через `docker-compose.yml`.

* `make clean`
  Удалить артефакты сборки.

---

## 🐳 Docker & Docker Compose

В корне находится `docker-compose.yml`, который поднимает:

* **setup**: генерирует сертификаты (CA и node-certs)
* **es01/es02/es03**: трёхнодовый кластер Elasticsearch (без HTTP-SSL, `xpack.security.enabled=false`)
* **kibana**: для просмотра индексов

```bash
# Поднять ES-кластер и Kibana
make docker-up

# Остановить
make docker-down
```

---

## 🧪 Переменные окружения

Создайте файл `.env` в корне с такими ключами:

```env
# Elasticsearch
ES_PORT=127.0.0.1:9200
ELASTIC_USER=elastic
ELASTIC_PASSWORD=123456

# ES cluster settings
STACK_VERSION=9.0.0
CLUSTER_NAME=docker-cluster
LICENSE=basic
MEM_LIMIT=1073741824

# HTTP-пул
ES_MAX_IDLE_CONNS=100
ES_MAX_IDLE_CONNS_PER_HOST=10
ES_IDLE_CONN_TIMEOUT=30

# (Опционально) Kibana
KIBANA_PORT=5601
```

---

## 🔌 API Endpoints

### POST `/product/search`

Поиск товаров
**Request** (`application/json`):

```json
{
  "query": "laptop",
  "categories": ["electronics", "computers"],
  "brand": ["Dell", "HP"],
  "min_price": 500,
  "max_price": 2000,
  "sort_by": "price",
  "sort_order": "asc",
  "page": 1,
  "page_size": 10,
  "highlight_fields": ["name","description"]
}
```

**Response** (`200 OK`, `application/json`):

```json
{
  "products":[
    {"id":"1","name":"Dell XPS 13","description":"...","price":999,"stock":42,"category":"electronics","brand":"Dell"},
    // ...
  ],
  "total": 42,
  "page": 1,
  "page_size": 10,
  "highlights": {"1":["<em>Dell</em> XPS 13"]},
  "facets": {
    "category":[{"key":"electronics","count":42}],
    "brand":[{"key":"Dell","count":23}]
  }
}
```

### GET `/product/health`

Проверка состояния кластера
**Response** (`200 OK`, `application/json`):

```json
{
  "status":"green",
  "active_shards":3,
  "relocating_shards":0,
  "unassigned_shards":0,
  "timed_out":false
}
```

---

## 📝 DTO

Все HTTP-DTO лежат в `transport/rest/product/product_dto`:

```go
// SearchRequest — параметры поиска
type SearchRequest struct {
  Query           string   `json:"query"`
  Categories      []string `json:"categories"`
  Brand           []string `json:"brand"`
  MinPrice        float64  `json:"min_price"`
  MaxPrice        float64  `json:"max_price"`
  SortBy          string   `json:"sort_by"`
  SortOrder       string   `json:"sort_order"`
  Page            int      `json:"page"`
  PageSize        int      `json:"page_size"`
  HighlightFields []string `json:"highlight_fields"`
}

// SearchResponse — ответ поиска
type SearchResponse struct {
  Products   []*entity.Product       `json:"products"`
  Total      int64                   `json:"total"`
  Page       int                     `json:"page"`
  PageSize   int                     `json:"page_size"`
  Highlights map[string][]string     `json:"highlights"`
  Facets     map[string]entity.Facet `json:"facets"`
}

// HealthResponse — ответ health
type HealthResponse struct {
  Status           string `json:"status"`
  ActiveShards     int64  `json:"active_shards"`
  RelocatingShards int64  `json:"relocating_shards"`
  UnassignedShards int64  `json:"unassigned_shards"`
  TimedOut         bool   `json:"timed_out"`
}
```

---

## 📚 Swagger Documentation

1. Установите `swag`:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Сгенерируйте docs:

   ```bash
   swag init -g cmd/search_service/main.go
   ```

3. Swagger UI доступен по адресу:

   ```
   http://localhost:8080/swagger/index.html
   ```

---

## 🙌 Логирование

Везде используется `log` с префиксами:

| Префикс                   | Слой                     |
| ------------------------- | ------------------------ |
| `[transport]`             | HTTP-хендлеры (gin)      |
| `[use_cases]`             | Бизнес-логика            |
| `[repository]`            | ES-репозитории           |
| `[kafka_product_service]` | Kafka consumer service   |
| `[elasticsearch]`         | Подключение и ping       |
| `[db]`                    | Инициализация соединений |

**Пример:**

```
[transport] SearchProducts called
[use_cases] SearchProducts succeeded: found 42 items
[repository] Bulk save succeeded for 5 products
[kafka_product_service][consumer-0] message processed successfully
```

---

С этим README вы быстро настроите, запустите и поймёте архитектуру и возможности **product-search**!
