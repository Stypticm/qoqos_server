import requests
import json
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import uvicorn

app = FastAPI()

# Конфигурация
OLLAMA_URL = "http://localhost:11434/api"
QDRANT_URL = "http://localhost:6333/collections/qoqos-knowledge/points/search"
EMBED_MODEL = "nomic-embed-text"
CHAT_MODEL = "llama3.2:3b"

class ChatRequest(BaseModel):
    message: string

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

        # 3. Формируем финальный ответ через LLM
        system_prompt = f"""Ты — Валера, экспертный ассистент системы QOQOS. Твоя задача — давать ОЧЕНЬ ТОЧНЫЕ ответы, основываясь ИСКЛЮЧИТЕЛЬНО на предоставленном КОНТЕКСТЕ.
ПРАВИЛА:
1. Если в КОНТЕКСТЕ есть стадии или списки, цитируй их точно.
2. Не придумывай информацию от себя.
3. Если данных нет — просто ответь: "В моей базе знаний нет информации об этом".
4. Отвечай на русском языке.

КОНТЕКСТ:
{context}
"""

        res = requests.post(f"{OLLAMA_URL}/generate", json={
            "model": CHAT_MODEL,
            "prompt": question,
            "system": system_prompt,
            "stream": False
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
    print("🚀 Brain Service (Valera) started on port 8001")
    uvicorn.run(app, host="0.0.0.0", port=8001)
