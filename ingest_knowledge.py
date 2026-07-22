import os
import sys
import requests
import uuid

OLLAMA_URL = "http://localhost:11434/api"
QDRANT_URL = "http://localhost:6333/collections/qoqos-knowledge/points"
EMBED_MODEL = "nomic-embed-text"

def get_embedding(text):
    res = requests.post(f"{OLLAMA_URL}/embeddings", json={
        "model": EMBED_MODEL,
        "prompt": text
    })
    res.raise_for_status()
    return res.json()["embedding"]

def ingest_file(filepath):
    print(f"Читаем файл {filepath}...")
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()
    
    # Разбивка на абзацы (чанки)
    chunks = [c.strip() for c in content.split("\n\n") if c.strip()]
    
    points = []
    for chunk in chunks:
        embedding = get_embedding(chunk)
        # Генерируем детерминированный ID на основе текста, чтобы избежать дубликатов при повторном запуске
        point_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, chunk.encode('utf-8').hex()))
        points.append({
            "id": point_id,
            "vector": embedding,
            "payload": {
                "text": chunk,
                "source": os.path.basename(filepath)
            }
        })
    
    if points:
        res = requests.put(QDRANT_URL, json={"points": points})
        res.raise_for_status()
        # Используем безопасный print, чтобы избежать ошибок кодировки в консоли Windows
        print(f"Успешно загружено {len(points)} фрагментов из {filepath} в Qdrant.")

if __name__ == "__main__":
    # Если передан конкретный файл: python ingest_knowledge.py QOQOS_KNOWLEDGE/my_rules.txt
    if len(sys.argv) > 1:
        filepath = sys.argv[1]
        if os.path.exists(filepath):
            ingest_file(filepath)
        else:
            print(f"Файл {filepath} не найден.")
    else:
        # Иначе загружаем все .txt файлы из папки
        knowledge_dir = "QOQOS_KNOWLEDGE"
        for filename in os.listdir(knowledge_dir):
            if filename.endswith(".txt"):
                filepath = os.path.join(knowledge_dir, filename)
                ingest_file(filepath)
