-- Таблица для email событий (outbox worker)
CREATE TABLE IF NOT EXISTS outbox_events (
    event_id     Int32,
    event_type   String,
    flat_id      Int32,
    house_id     Int32,
    email        String,
    status       String,       -- sent / failed
    created_at   DateTime,
    processed_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (processed_at, event_type);

-- Таблица для всех бизнес-событий (аналитика)
CREATE TABLE IF NOT EXISTS business_events (
    event_id     UUID,
    event_type   String,         -- 'user_registered', 'house_created', 'flat_created', etc.
    user_id      String,         -- UUID пользователя (если известен)
    entity_id    Int32,          -- ID дома/квартиры (если применимо)
    entity_type  String,         -- 'house', 'flat', 'user', 'subscription'
    metadata     String,         -- JSON с дополнительными данными
    status       String,         -- 'success', 'failed'
    error        String,         -- Текст ошибки (если failed)
    timestamp    DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (timestamp, event_type)
PARTITION BY toYYYYMM(timestamp);