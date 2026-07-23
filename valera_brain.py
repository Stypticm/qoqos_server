import requests
import json
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# Конфигурация
OLLAMA_URL = "http://localhost:11434/api"
QDRANT_URL = "http://localhost:6333/collections/qoqos-knowledge/points/search"
EMBED_MODEL = "nomic-embed-text"
CHAT_MODEL = "qwen2:1.5b"

class ChatRequest(BaseModel):
    message: str

def get_brain_answer(question):
    try:
        # 1. Получаем эмбеддинг вопроса
        res = requests.post(f"{OLLAMA_URL}/embeddings", json={
            "model": EMBED_MODEL,
            "prompt": question
        })
        res.raise_for_status()
        embedding = res.json()["embedding"]

        # 2. Ищем контекст в Qdrant
        res = requests.post(QDRANT_URL, json={
            "vector": embedding,
            "limit": 7,
            "with_payload": True
        })
        res.raise_for_status()
        hits = res.json()["result"]
        
        context = ""
        for hit in hits:
            if hit.get("payload") and hit["payload"].get("text"):
                source = hit['payload'].get('source', 'базы')
                text = hit['payload']['text']
                context += f"\n--- Фрагмент из {source}: ---\n{text}\n"

        if len(context) > 5000:
            context = context[:5000] + "... [Контекст обрезан]"

        if not context.strip():
            context = "Информации в базе знаний не найдено."

        if context.strip() == "Информации в базе знаний не найдено.":
            return "<ESCALATE_REASON>Вопрос не найден в базе знаний</ESCALATE_REASON>"

        # 3. Формируем финальный ответ через LLM
        system_prompt = f"""Ты — Валера, помощник сервиса QOQOS по выкупу оборудования. Отвечай на русском языке СТРОГО на основе информации ниже.

ПРАВИЛА:
- Если в контексте есть информация, ответь на вопрос клиента.
- Обязательно проверь, не отвечает ли контекст на вопрос про QOQOS, выкуп (trade-in), оценку, устройство, услугу — даже если ответ частичный, дай его.
- ТОЛЬКО если в контексте нет НИЧЕГО релевантного вопросу, выведи ПЕРВОЙ строкой БЕЗ ЛЮБОГО ДОПОЛНЕНИЯ: <ESCALATE_REASON>Вопрос про ""{question[:60]}"" не найден в базе знаний</ESCALATE_REASON> и затем напиши отдельной строкой, что оператор поможет уточнить. Никогда не изменяй этот тег, не добавляй к нему префиксы, не заключай в кавычки, не заменяй угловые скобки.
- Если контекст содержит хоть какую-то релевантную информацию — отвечай сразу, без тегов.

ФОРМАТ ДЛЯ ESCALATE:
<ESCALATE_REASON>Вопрос про "..." не найден в базе знаний</ESCALATE_REASON>
Оператор поможет уточнить информацию.

ИНФОРМАЦИЯ ДЛЯ ОТВЕТА:
{context}
"""

        full_prompt = f"{system_prompt}\nВопрос клиента: {question}\nОтвет Валеры:"

        res = requests.post(f"{OLLAMA_URL}/generate", json={
            "model": CHAT_MODEL,
            "prompt": full_prompt,
            "system": "",
            "stream": False,
            "options": {
                "temperature": 0.3,
                "num_predict": 300,
                "stop": ["Вопрос клиента:", "\n\n"]
            }
        })
        res.raise_for_status()
        return res.json().get("response", "Ошибка генерации")
    except Exception as e:
        print(f"Error in brain: {e}")
        return "Извините, произошла техническая ошибка при поиске в базе знаний."

@app.post("/ask")
async def ask_valera(request: ChatRequest):
    answer = get_brain_answer(request.message)
    return {"response": answer}

if __name__ == "__main__":
    print("Brain Service (Valera) started on port 8001")
    uvicorn.run(app, host="0.0.0.0", port=8001)
