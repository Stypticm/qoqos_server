import requests

def create_collection():
    url = "http://localhost:6333/collections/qoqos-knowledge"
    payload = {
        "vectors": {
            "size": 768,
            "distance": "Cosine"
        }
    }
    try:
        r = requests.put(url, json=payload)
        r.raise_for_status()
        print(f"✅ Результат: {r.json()}")
    except Exception as e:
        print(f"❌ Ошибка: {e}")

if __name__ == "__main__":
    create_collection()
