import psycopg2
import json
import os

# Данные из .env
DB_USER = "qoqos"
DB_PASS = "StypticQtweDB"
DB_NAME = "qoqos_db"
DB_HOST = "localhost"
DB_PORT = "5432"

try:
    conn = psycopg2.connect(
        dbname=DB_NAME,
        user=DB_USER,
        password=DB_PASS,
        host=DB_HOST,
        port=DB_PORT
    )
    cur = conn.cursor()
    
    # 1. Достаем Блог
    cur.execute('SELECT * FROM "BlogPost";')
    cols = [desc[0] for desc in cur.description]
    blog_posts = [dict(zip(cols, row)) for row in cur.fetchall()]
    
    # 2. Достаем Лоты (если нужно для каталога)
    cur.execute('SELECT * FROM "MarketplaceLot";')
    cols = [desc[0] for desc in cur.description]
    lots = [dict(zip(cols, row)) for row in cur.fetchall()]
    
    data = {
        "blog_posts": blog_posts,
        "lots": lots
    }
    
    with open('db_data.json', 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2, default=str)
        
    print("Successfully exported DB data to db_data.json")

except Exception as e:
    print(f"Error: {e}")
finally:
    if 'conn' in locals():
        conn.close()
