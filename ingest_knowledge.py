#!/usr/bin/env python3
"""
Ingest knowledge base files into ai-agent data files.

Usage:
    python ingest_knowledge.py

Reads:
    QOQOS_KNOWLEDGE/defects.txt       -> defects/questions for Qdrant
    QOQOS_KNOWLEDGE/repair_prices.csv -> prices for Postgres
    QOQOS_KNOWLEDGE/stock.csv         -> stock for Postgres

Updates:
    ai-agent/load_qdrant.py
    migrations/seed_repair_data.sql
"""

import os
import re
import csv
import json
import asyncio
import httpx
from pathlib import Path
from typing import List, Dict, Optional

# === CONFIG ===
BASE_DIR = Path(__file__).parent
QOQOS_KNOWLEDGE_DIR = BASE_DIR / "QOQOS_KNOWLEDGE"
AI_AGENT_DIR = BASE_DIR / "ai-agent"
MIGRATIONS_DIR = BASE_DIR / "migrations"

OLLAMA_EMBEDDING_URL = os.getenv("OLLAMA_EMBEDDING_URL", "http://localhost:11434/api/embeddings")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "nomic-embed-text")

# === HELPERS ===

def split_into_blocks(text: str) -> List[str]:
    lines = text.splitlines()
    blocks = []
    current = []
    for line in lines:
        stripped = line.strip()
        if stripped.startswith("#"):
            if current:
                blocks.append(" ".join(current).strip())
                current = []
            continue
        if stripped:
            current.append(stripped)
        else:
            if current:
                blocks.append(" ".join(current).strip())
                current = []
    if current:
        blocks.append(" ".join(current).strip())
    return [b for b in blocks if b]


def normalize_part_name(text: str) -> tuple[Optional[str], Optional[str]]:
    text_lower = text.lower()
    part_aliases = {
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
    for alias, canonical in part_aliases.items():
        if alias in text_lower:
            return canonical, None
    return None, None


def detect_model(text: str) -> Optional[str]:
    text_lower = text.lower()
    patterns = [
        r"iphone\s*\d+[a-z]?(?:\s*pro(?:\s*max)?)?(?:\s*plus)?(?:\s*mini)?",
        r"samsung\s*a\d+[a-z]?(?:\s*plus)?(?:\s*ultra)?",
        r"ipad\s*\d+[a-z]?(?:\s*pro)?(?:\s*mini)?",
        r"macbook\s*(?:air|pro)?\s*\d+[\"']?(?:\s*m\d+)?",
        r"airpods\s*\d+",
        r"redmi\s*\d+[a-z]?",
        r"xiaomi\s*\d+[a-z]?",
    ]
    for pattern in patterns:
        match = re.search(pattern, text_lower)
        if match:
            return match.group(0).strip().title()
    return None


async def get_embedding(text: str) -> List[float]:
    async with httpx.AsyncClient() as client:
        response = await client.post(
            OLLAMA_EMBEDDING_URL,
            json={"model": EMBEDDING_MODEL, "prompt": text},
            timeout=60.0,
        )
        response.raise_for_status()
        return response.json()["embedding"]


# === LOADERS ===

def load_txt_knowledge() -> List[Dict]:
    defects = []
    txt_file = QOQOS_KNOWLEDGE_DIR / "defects.txt"
    if not txt_file.exists():
        return defects

    text = txt_file.read_text(encoding="utf-8")
    blocks = split_into_blocks(text)
    for block in blocks:
        if len(block) < 10:
            continue
        part, _ = normalize_part_name(block)
        model = detect_model(block)
        defects.append({
            "text": block,
            "part": part,
            "model": model,
        })
    return defects


def load_csv_knowledge() -> tuple[List[Dict], List[Dict]]:
    prices = []
    stock = []

    for csv_file in QOQOS_KNOWLEDGE_DIR.glob("*.csv"):
        with open(csv_file, "r", encoding="utf-8") as f:
            reader = csv.DictReader(f)
            rows = list(reader)

        if csv_file.name == "stock.csv":
            for row in rows:
                stock.append({
                    "part_name": row.get("part_name", "").strip(),
                    "model": row.get("model", "").strip(),
                    "price": int(row.get("price", 0)),
                    "quantity": int(row.get("quantity", 0)),
                    "delivery_days": int(row.get("delivery_days", 3)) if row.get("delivery_days") else None,
                })
        else:
            for row in rows:
                prices.append({
                    "part_name": row.get("part_name", "").strip(),
                    "model": row.get("model", "").strip(),
                    "price": int(row.get("price", 0)),
                    "quantity": int(row.get("quantity", 0)),
                    "delivery_days": int(row.get("delivery_days", 3)) if row.get("delivery_days") else None,
                })

    return prices, stock


# === FILE GENERATORS ===

def generate_load_qdrant(defects: List[Dict]) -> str:
    lines = [
        "import os",
        "import asyncio",
        "import httpx",
        "from qdrant_client import QdrantClient",
        "from qdrant_client.models import PointStruct, VectorParams, Distance",
        "",
        "QDRANT_HOST = os.getenv('QDRANT_HOST', 'localhost')",
        "QDRANT_PORT = int(os.getenv('QDRANT_PORT', '6333'))",
        "QDRANT_COLLECTION = os.getenv('QDRANT_COLLECTION', 'defects')",
        "OLLAMA_EMBEDDING_URL = os.getenv('OLLAMA_EMBEDDING_URL', 'http://localhost:11434/api/embeddings')",
        "EMBEDDING_MODEL = os.getenv('EMBEDDING_MODEL', 'nomic-embed-text')",
        "",
        "qdrant = QdrantClient(host=QDRANT_HOST, port=QDRANT_PORT)",
        "",
        "defects_data = [",
    ]

    for item in defects:
        part_json = json.dumps(item["part"], ensure_ascii=False)
        model_json = json.dumps(item["model"], ensure_ascii=False)
        lines.append(f'    {{"text": {json.dumps(item["text"], ensure_ascii=False)}, "part": {part_json}, "model": {model_json}}},')

    lines.extend([
        "]",
        "",
        "async def get_embedding(text: str) -> list[float]:",
        "    async with httpx.AsyncClient() as client:",
        "        response = await client.post(",
        "            OLLAMA_EMBEDDING_URL,",
        "            json={'model': EMBEDDING_MODEL, 'prompt': text},",
        "            timeout=60.0,",
        "        )",
        "        response.raise_for_status()",
        "        return response.json()['embedding']",
        "",
        "async def main():",
        "    if not qdrant.collection_exists(QDRANT_COLLECTION):",
        "        qdrant.create_collection(",
        "            collection_name=QDRANT_COLLECTION,",
        "            vectors_config=VectorParams(size=768, distance=Distance.COSINE),",
        "        )",
        "",
        "    points = []",
        "    for i, item in enumerate(defects_data):",
        "        embedding = await get_embedding(item['text'])",
        "        points.append(",
        "            PointStruct(",
        "                id=i,",
        "                vector=embedding,",
        "                payload={",
        "                    'text': item['text'],",
        "                    'part': item['part'],",
        "                    'model': item['model'],",
        "                },",
        "            )",
        "        )",
        "",
        "    qdrant.upsert(collection_name=QDRANT_COLLECTION, points=points)",
        "    print(f'Loaded {len(points)} defects into Qdrant')",
        "",
        "if __name__ == '__main__':",
        "    asyncio.run(main())",
    ])

    return "\n".join(lines)


def generate_seed_sql(prices: List[Dict], stock: List[Dict]) -> str:
    lines = [
        "CREATE TABLE IF NOT EXISTS repair_prices (",
        "    id SERIAL PRIMARY KEY,",
        "    part_name VARCHAR(100) NOT NULL,",
        "    model VARCHAR(50) NOT NULL,",
        "    price INTEGER NOT NULL,",
        "    created_at TIMESTAMP DEFAULT NOW(),",
        "    updated_at TIMESTAMP DEFAULT NOW()",
        ");",
        "",
        "CREATE TABLE IF NOT EXISTS stock (",
        "    id SERIAL PRIMARY KEY,",
        "    part_name VARCHAR(100) NOT NULL,",
        "    model VARCHAR(50) NOT NULL,",
        "    quantity INTEGER NOT NULL DEFAULT 0,",
        "    delivery_days INTEGER DEFAULT NULL,",
        "    created_at TIMESTAMP DEFAULT NOW(),",
        "    updated_at TIMESTAMP DEFAULT NOW()",
        ");",
        "",
        "CREATE INDEX IF NOT EXISTS idx_repair_prices_part_model ON repair_prices(part_name, model);",
        "CREATE INDEX IF NOT EXISTS idx_stock_part_model ON stock(part_name, model);",
        "",
    ]

    if prices:
        lines.append("INSERT INTO repair_prices (part_name, model, price) VALUES")
        price_lines = []
        for p in prices:
            price_lines.append(f"('{p['part_name']}', '{p['model']}', {p['price']})")
        lines.append(",\n".join(price_lines))
        lines.append("ON CONFLICT DO NOTHING;")
        lines.append("")
    else:
        lines.extend([
            "INSERT INTO repair_prices (part_name, model, price) VALUES",
            "('экран', 'iPhone 14', 4990),",
            "('экран', 'iPhone 13', 3990),",
            "('экран', 'Samsung S23', 3500),",
            "('кнопка громкости', 'Samsung A52', 1200),",
            "('аккумулятор', 'iPhone 14', 3500)",
            "ON CONFLICT DO NOTHING;",
            "",
        ])

    if stock:
        lines.append("INSERT INTO stock (part_name, model, quantity, delivery_days) VALUES")
        stock_lines = []
        for s in stock:
            delivery = s["delivery_days"] if s["delivery_days"] is not None else "NULL"
            stock_lines.append(f"('{s['part_name']}', '{s['model']}', {s['quantity']}, {delivery})")
        lines.append(",\n".join(stock_lines))
        lines.append("ON CONFLICT DO NOTHING;")
    else:
        lines.extend([
            "INSERT INTO stock (part_name, model, quantity, delivery_days) VALUES",
            "('экран', 'iPhone 14', 5, 0),",
            "('экран', 'iPhone 13', 0, 3),",
            "('кнопка громкости', 'Samsung A52', 0, 5)",
            "ON CONFLICT DO NOTHING;",
        ])

    lines.append("")
    return "\n".join(lines)


# === MAIN ===

def main():
    print("Scanning QOQOS_KNOWLEDGE...")
    defects = load_txt_knowledge()
    prices, stock = load_csv_knowledge()

    print(f"   Found {len(defects)} defect blocks")
    print(f"   Found {len(prices)} price rows")
    print(f"   Found {len(stock)} stock rows")

    # Update load_qdrant.py
    qdrant_file = AI_AGENT_DIR / "load_qdrant.py"
    qdrant_code = generate_load_qdrant(defects)
    qdrant_file.write_text(qdrant_code, encoding="utf-8")
    print(f"Updated {qdrant_file}")

    # Update seed SQL
    seed_file = MIGRATIONS_DIR / "seed_repair_data.sql"
    seed_sql = generate_seed_sql(prices, stock)
    seed_file.write_text(seed_sql, encoding="utf-8")
    print(f"Updated {seed_file}")

    print("\nNext steps:")
    print("  1. Review generated files")
    print("  2. Reload Qdrant:  docker exec qoqos_ai_agent python /app/load_qdrant.py")
    print("  3. Reload Postgres:")
    print("       docker cp qoqos_server/migrations/seed_repair_data.sql qoqos_postgres:/tmp/seed.sql")
    print("       docker exec qoqos_postgres psql -U qoqos -d qoqos_db -f /tmp/seed.sql")


if __name__ == "__main__":
    main()
