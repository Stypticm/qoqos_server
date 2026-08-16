import asyncpg
import json
import os
import re
from qdrant_client import QdrantClient
import httpx
import pymorphy3

# === КОНФИГУРАЦИЯ ===
DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = os.getenv("DB_PORT", "5432")
DB_USER = os.getenv("DB_USER", "postgres")
DB_PASSWORD = os.getenv("DB_PASSWORD", "postgres")
DB_NAME = os.getenv("DB_NAME", "qoqos")
DB_DSN = f"postgresql://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"

QDRANT_HOST = os.getenv("QDRANT_HOST", "localhost")
QDRANT_PORT = int(os.getenv("QDRANT_PORT", "6333"))
QDRANT_COLLECTION = os.getenv("QDRANT_COLLECTION", "defects")

OLLAMA_URL = os.getenv("OLLAMA_URL", "http://localhost:11434/api/generate")
OLLAMA_EMBEDDING_URL = os.getenv("OLLAMA_EMBEDDING_URL", "http://localhost:11434/api/embeddings")
OLLAMA_MODEL = os.getenv("OLLAMA_MODEL", "deepseek-r1:8b")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "nomic-embed-text")

# === КЛИЕНТЫ ===
qdrant = QdrantClient(host=QDRANT_HOST, port=QDRANT_PORT)
morph = pymorphy3.MorphAnalyzer()

PART_ALIASES = {
    "экран": "экран",
    "дисплей": "экран",
    "матрица": "экран",
    "стекло": "экран",
    "кнопка": "кнопка громкости",
    "клавиша": "кнопка громкости",
    "громкость": "кнопка громкости",
    "аккумулятор": "аккумулятор",
    "батарея": "аккумулятор",
    "батарейка": "аккумулятор",
    "зарядка": "аккумулятор",
    "разъем": "разъем зарядки",
    "гнездо": "разъем зарядки",
}

def normalize_part_name(user_part: str) -> str:
    user_part_lower = user_part.lower().strip()
    if user_part_lower in PART_ALIASES:
        return PART_ALIASES[user_part_lower]
    for alias, canonical in PART_ALIASES.items():
        if alias in user_part_lower:
            return canonical
    return user_part

async def get_embedding(text: str) -> list[float]:
    try:
        async with httpx.AsyncClient() as client:
            response = await client.post(
                OLLAMA_EMBEDDING_URL,
                json={"model": EMBEDDING_MODEL, "prompt": text},
                timeout=60.0,
            )
            response.raise_for_status()
            return response.json()["embedding"]
    except Exception as e:
        print(f"Ошибка получения эмбеддинга: {e}")
        return []

async def get_price(part_name: str, model: str):
    try:
        conn = await asyncpg.connect(DB_DSN)
        row = await conn.fetchrow(
            "SELECT price FROM repair_prices WHERE part_name = $1 AND model = $2",
            part_name, model,
        )
        await conn.close()
        return row['price'] if row else None
    except Exception as e:
        print(f"Ошибка get_price: {e}")
        return None

async def check_stock(part_name: str, model: str):
    try:
        conn = await asyncpg.connect(DB_DSN)
        row = await conn.fetchrow(
            "SELECT quantity, delivery_days FROM stock WHERE part_name = $1 AND model = $2",
            part_name, model,
        )
        await conn.close()
        if not row:
            return {"in_stock": False, "delivery_days": None}
        return {"in_stock": row['quantity'] > 0, "delivery_days": row['delivery_days']}
    except Exception as e:
        print(f"Ошибка check_stock: {e}")
        return {"in_stock": False, "delivery_days": None}

async def search_defects(query: str, top_k: int = 3):
    try:
        embedding = await get_embedding(query)
        if not embedding:
            return []
        results = qdrant.query_points(
            collection_name=QDRANT_COLLECTION,
            query=embedding,
            limit=top_k,
        )
        return [{
            "text": hit.payload["text"],
            "part": hit.payload.get("part"),
            "model": hit.payload.get("model"),
            "score": hit.score,
        } for hit in results.points]
    except Exception as e:
        print(f"Ошибка search_defects: {e}")
        return []

async def call_ollama(prompt: str) -> str:
    try:
        async with httpx.AsyncClient() as client:
            response = await client.post(
                OLLAMA_URL,
                json={
                    "model": OLLAMA_MODEL,
                    "prompt": prompt,
                    "stream": False,
                    "temperature": 0.3,
                },
                timeout=60.0,
            )
            response.raise_for_status()
            return response.json()["response"]
    except Exception as e:
        print(f"Ошибка Ollama: {e}")
        return json.dumps({"part_name": None, "model": None, "message": "Ошибка вызова ИИ"})

async def process_query(user_query: str) -> str:
    try:
        similar_defects = await search_defects(user_query)
        if not similar_defects:
            context_text = "Похожих дефектов не найдено."
        else:
            context_text = "\n".join([f"- {d['text']}" for d in similar_defects])

        prompt = f"""
Ты — AI-агент сервиса по ремонту телефонов.

Вот примеры правильных ответов:

Пример 1:
Вопрос пользователя: "Сколько стоит замена экран на iPhone 14?"
Правильный JSON: {{"part_name": "экран", "model": "iPhone 14", "message": ""}}

Пример 2:
Вопрос пользователя: "У меня сломался экран на айфоне 14, сколько?"
Правильный JSON: {{"part_name": "экран", "model": "iPhone 14", "message": ""}}

Пример 3:
Вопрос пользователя: "Замена кнопку громкости на самсунг а52"
Правильный JSON: {{"part_name": "кнопка громкости", "model": "Samsung A52", "message": ""}}

Пример 4:
Вопрос пользователя: "Телефон быстро разряжается, нужна замена аккумулятора"
Правильный JSON: {{"part_name": "аккумулятор", "model": "неизвестно", "message": "Уточните модель телефона"}}

Теперь твоя очередь. Ответь на вопрос пользователя строго в формате JSON, как в примерах выше.
Не используй лишние слова, только JSON.

Вопрос пользователя: "{user_query}"
Ответ:
"""

        response = await call_ollama(prompt)
        print("Ollama ответил:", response)

        clean_response = response.strip()
        if clean_response.startswith("```json"):
            clean_response = clean_response[7:]
        if clean_response.endswith("```"):
            clean_response = clean_response[:-3]
        clean_response = clean_response.strip()

        clean_response = clean_response.replace("part_ name", "part_name")
        clean_response = re.sub(r'(?<!")(\bmodel\b)(?!")(\s*:)', r'"\1"\2', clean_response)
        clean_response = re.sub(r'(?<!")(\bpart_name\b)(?!")(\s*:)', r'"\1"\2', clean_response)

        try:
            data = json.loads(clean_response)
        except json.JSONDecodeError as e:
            print(f"Ошибка парсинга JSON: {e}")
            print(f"Строка: {clean_response}")
            return "Извините, не удалось обработать запрос. Попробуйте уточнить: 'Экран iPhone 14'."

        if data.get("message"):
            return data["message"]

        part = data.get("part_name")
        model = data.get("model")

        if not part or not model or part == "null" or model == "null":
            return "Уточните, пожалуйста, модель телефона и что именно сломалось (экран, кнопка, аккумулятор)."

        price = await get_price(part, model)
        stock = await check_stock(part, model)

        if price is None:
            return f"К сожалению, мы не нашли цену для {part} на {model}. Пожалуйста, уточните запрос."

        in_stock = False
        days = 3
        
        if stock:
            in_stock = stock.get("in_stock", False)
            days = stock.get("delivery_days")
            if days is None:
                days = 3

        if in_stock:
            return f"✅ Замена {part} на {model} стоит {price} ₽. Запчасть есть в наличии. Можем записать вас на ремонт завтра. Вас устраивает?"
        else:
            return f"⚠️ Замена {part} на {model} стоит {price} ₽. К сожалению, запчасть отсутствует, но мы можем заказать её (срок поставки {days} дня). Хотите оставить заявку?"

    except Exception as e:
        print(f"Критическая ошибка в process_query: {e}")
        return "Произошла внутренняя ошибка. Попробуйте ещё раз или уточните запрос."

def normalize_part_name(user_part: str) -> str:
    normalized = morph.parse(user_part)[0].normal_form
    return normalized