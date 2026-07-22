import requests
import json

def debug_qdrant_data():
    url = "http://localhost:6333/collections/qoqos-knowledge/points/scroll"
    try:
        r = requests.post(url, json={"limit": 5, "with_payload": True})
        r.raise_for_status()
        points = r.json()["result"]["points"]
        print(f"📊 Всего точек (на проверку): {len(points)}")
        for i, pt in enumerate(points):
            payload = pt.get("payload", {})
            print(f"\n--- Точка {i+1} (ID: {pt['id']}) ---")
            print(f"Payload keys: {list(payload.keys())}")
            for key, value in payload.items():
                print(f"[{key}]: {str(value)[:200]}...")
    except Exception as e:
        print(f"Ошибка: {e}")

if __name__ == "__main__":
    debug_qdrant_data()
