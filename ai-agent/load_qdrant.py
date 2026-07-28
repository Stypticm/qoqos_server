import os
import asyncio
import httpx
from qdrant_client import QdrantClient
from qdrant_client.models import PointStruct, VectorParams, Distance

QDRANT_HOST = os.getenv('QDRANT_HOST', 'localhost')
QDRANT_PORT = int(os.getenv('QDRANT_PORT', '6333'))
QDRANT_COLLECTION = os.getenv('QDRANT_COLLECTION', 'defects')
OLLAMA_EMBEDDING_URL = os.getenv('OLLAMA_EMBEDDING_URL', 'http://localhost:11434/api/embeddings')
EMBEDDING_MODEL = os.getenv('EMBEDDING_MODEL', 'nomic-embed-text')

qdrant = QdrantClient(host=QDRANT_HOST, port=QDRANT_PORT)

defects_data = [
    {"text": "Треснул экран iPhone 14, требуется замена дисплея", "part": "экран", "model": "Iphone 14"},
    {"text": "Желтизна экрана после падения, нужна замена матрицы", "part": "экран", "model": null},
    {"text": "Пожелтение экрана от времени, естественный износ", "part": "экран", "model": null},
    {"text": "Заводской брак дисплея на iPhone 14, ремонт по гарантии", "part": null, "model": "Iphone 14"},
    {"text": "Не работает кнопка громкости на Samsung A52", "part": "кнопка громкости", "model": "Samsung A52"},
    {"text": "Телефон быстро разряжается, нужна замена аккумулятора iPhone 14", "part": "аккумулятор", "model": "Iphone 14"},
    {"text": "Как работает выкуп (Trade-in) в QOQOS", "part": null, "model": null},
    {"text": "Сколько стоит замена экрана на iPhone 14", "part": "экран", "model": "Iphone 14"},
    {"text": "Нужна замена аккумулятора на Samsung", "part": "аккумулятор", "model": null},
    {"text": "Не работает разъем зарядки на iPhone", "part": "разъем зарядки", "model": null},
    {"text": "Как долго длится ремонт экрана", "part": "экран", "model": null},
    {"text": "Есть ли гарантия на ремонт", "part": null, "model": null},
    {"text": "Можно ли сдать телефон с разбитым экраном", "part": "экран", "model": null},
    {"text": "Сколько времени занимает диагностика", "part": null, "model": null},
    {"text": "Принимаете ли вы технику с дефектами", "part": null, "model": null},
    {"text": "Нужна ли предварительная запись на ремонт", "part": null, "model": null},
]

async def get_embedding(text: str) -> list[float]:
    async with httpx.AsyncClient() as client:
        response = await client.post(
            OLLAMA_EMBEDDING_URL,
            json={'model': EMBEDDING_MODEL, 'prompt': text},
            timeout=60.0,
        )
        response.raise_for_status()
        return response.json()['embedding']

async def main():
    if not qdrant.collection_exists(QDRANT_COLLECTION):
        qdrant.create_collection(
            collection_name=QDRANT_COLLECTION,
            vectors_config=VectorParams(size=768, distance=Distance.COSINE),
        )

    points = []
    for i, item in enumerate(defects_data):
        embedding = await get_embedding(item['text'])
        points.append(
            PointStruct(
                id=i,
                vector=embedding,
                payload={
                    'text': item['text'],
                    'part': item['part'],
                    'model': item['model'],
                },
            )
        )

    qdrant.upsert(collection_name=QDRANT_COLLECTION, points=points)
    print(f'Loaded {len(points)} defects into Qdrant')

if __name__ == '__main__':
    asyncio.run(main())