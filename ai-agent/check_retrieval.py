import asyncio, httpx
from qdrant_client import QdrantClient

async def main():
    q = "Есть ли гарантия на ремонт?"
    r = await httpx.AsyncClient().post(
        "http://ollama:11434/api/embeddings",
        json={"model": "nomic-embed-text", "prompt": q},
        timeout=60,
    )
    emb = r.json()["embedding"]
    print("dim:", len(emb))
    c = QdrantClient(host="qdrant", port=6333)
    res = c.query_points(collection_name="defects", query=emb, limit=20)
    for p in res.points:
        print(round(p.score, 3), p.payload.get("text"))

asyncio.run(main())