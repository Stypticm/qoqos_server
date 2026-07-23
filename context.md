# Контекст проекта QOQOS

Справочник для разработки. Здесь — только то, что не покрыт README и не очевидно из кода.

## 🧠 Архитектура «Валеры» (valera_brain.py)

```
Пользователь → chat.post.ts → Valera (port 8001)
                                    │
                         ┌──────────┴──────────┐
                         │                     │
                    Qdrant (поиск        Ollama (LLM)
                    релевантных          + эмбеддинги)
                    чанков)
                         │                     │
                         └──────────┬──────────┘
                                    ▼
                          Ответ с <ESCALATE_REASON>
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
               Есть ответ                    Нет инфы в базе
                    │                               │
                    ▼                               ▼
            Показываем пользователю       Telegram → оператору
                                            │
                                            ▼
                                   Оператор отвечает
                                            │
                                            ▼
                                   chat_reply в Postgres
                                            │
                                            ▼
                                   Фронтенд poll'ит
                                   и показывает клиенту
```

- **Qdrant**: коллекция `qoqos-knowledge`, 768d, Cosine. Хранит чанки из `QOQOS_KNOWLEDGE/*.txt`.
- **Escalation**: если в ответе есть `<ESCALATE_REASON>` или паттерны типа `нет данных`, `оператор поможет` — `chat.post.ts` отправляет уведомление в Telegram (`TELEGRAM_CHAT_ID_misha`) и возвращает клиенту заглушку.
- **Polling**: фронтенд каждые 3с дергает `/api/agents/poll?guestId=...` → `telegram_bot.js` (port 8002) → Postgres `ChatReply`.

## 📂 База знаний (Qdrant)

- Файлы: `QOQOS_KNOWLEDGE/*.txt`
- Загрузка: `py ingest_knowledge.py` (чанки → эмбеддинги → upsert в Qdrant)
- При обновлении файлов — перезапустить `py ingest_knowledge.py`
- Отладка: `py debug_rag.py` (scroll точек)

## 📝 Правила разработки

- Git Flow: `dev` → `feature/name` → merge в `dev` → merge в `main`
- Секреты только в `.env`
- Команды линтинга/тестов смотреть в `AGENTS.md`

## 🛠 Быстрый старт

```powershell
# 1. Инфра
docker-compose up -d

# 2. Знания в Qdrant
py ingest_knowledge.py

# 3. Бот + мозг
npm run dev
```
