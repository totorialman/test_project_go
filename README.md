# Серверное приложение для мониторинга состояния сети и обнаружения потенциальных угроз

Проект реализован на языке **Go** с соблюдением принципов **Clean Architecture**, что обеспечивает высокую тестируемость, поддерживаемость и возможность легкой замены компонентов (БД, AI-модуль).

## Особенности

*   **Высокая производительность:** Обработка потоков сетевых логов (JSON lines в ZIP-архивах) с использованием горутин.
*   **Гибридное хранение данных:**
    *   **PostgreSQL:** Метаданные, пользователи, токены доступа, конфигурация.
    *   **ClickHouse:** Хранение и быстрый анализ больших объемов сетевых логов (TimeSeries).
*   **Безопасность:** Статическая токенная авторизация для агентов, хеширование чувствительных данных, поддержка HTTPS.
*   **AI-анализ:** Интеграция модуля машинного обучения для детектирования аномалий (сканирование портов, DDoS-паттерны).
*   **Уведомления:** Интеграция с Telegram Bot API для алертинга.
*   **Контейнеризация:** Полный стек развертывается через `docker-compose`.

## Технологический стек

| Компонент | Технология |
| :--- | :--- |
| **Язык** | Go 1.21+ |
| **Архитектура** | Clean Architecture (Layers: Handler, Service, Repository, Model) |
| **Web Framework** | gorilla/mux |
| **БД (Метаданные)** | PostgreSQL 15+ |
| **БД (Логи/Аналитика)** | ClickHouse 23.3+ |
| **Контейнеризация** | Docker, Docker Compose |
| **Testing** | `testing`, `testify`, `gomock` |

## API Documentation

API построено по принципам REST. Все запросы от агентов требуют наличия заголовка авторизации.

### Авторизация
Все защищенные эндпоинты требуют заголовок:
`Authorization: Bearer <your_static_token>`

Токен генерируется администратором при деплое или через админ-панель.

### Эндпоинты

#### 1. Загрузка логов 
Принимает ZIP-архив, содержащий файлы с JSON логами (формат JSON Lines).
*   **URL:** `POST /api/v1/logs/upload`
*   **Auth:** Required (Agent Token)
*   **Content-Type:** `multipart/form-data`
*   **Body:**
    *   `file`: ZIP архив.
*   **Response:** `200 OK`
    ```json
    {
      "status": "success",
      "message": "Logs accepted for processing",
      "batch_id": "uuid-v4-string"
    }
    ```

#### 2. Получение списка аномалий
Возвращает список выявленных инцидентов с возможностью фильтрации.
*   **URL:** `GET /api/v1/anomalies`
*   **Auth:** Required (Admin Token)
*   **Query Params:**
    *   `limit` (int, default: 20)
    *   `offset` (int, default: 0)
    *   `severity` (string: "low", "medium", "high", "critical")
    *   `from` (timestamp)
    *   `to` (timestamp)
*   **Response:** `200 OK`
    ```json
    {
      "total": 150,
      "items": [
        {
          "id": "uuid",
          "timestamp": "2026-05-20T10:00:00Z",
          "src_ip": "192.168.1.50",
          "dst_ip": "10.0.0.1",
          "type": "port_scan",
          "severity": "high",
          "status": "new",
          "details": { "ports_scanned": 1024, "duration_ms": 500 }
        }
      ]
    }
    ```

#### 3. Изменение статуса инцидента
Обновляет статус обработки инцидента (например, `new` -> `resolved`).
*   **URL:** `PUT /api/v1/anomalies/{id}`
*   **Auth:** Required (Admin Token)
*   **Body:**
    ```json
    {
      "status": "resolved",
      "comment": "False positive, internal scan"
    }
    ```
*   **Response:** `200 OK`

#### 4. Статистика 
Агрегированные данные для дашборда.
*   **URL:** `GET /api/v1/stats/summary`
*   **Auth:** Required (Admin Token)
*   **Query Params:** `period` (hour, day, week)
*   **Response:** `200 OK`
    ```json
    {
      "total_requests": 150000,
      "anomalies_detected": 12,
      "top_attacker_ips": ["192.168.1.50", "10.20.30.40"],
      "protocols_distribution": { "TCP": 80%, "UDP": 15%, "ICMP": 5% }
    }
    ```

#### 5. Управление токенами 
Генерация нового токена для агента.
*   **URL:** `POST /api/admin/tokens`
*   **Auth:** Required (Super Admin Token)
*   **Body:**
    ```json
    {
      "agent_name": "office-gateway-01",
      "description": "Main gateway logs"
    }
    ```
*   **Response:** `201 Created`
    ```json
    {
      "token": "static_secure_token_string",
      "created_at": "2026-05-20T12:00:00Z"
    }
    ```
