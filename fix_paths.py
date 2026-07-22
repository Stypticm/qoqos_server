import psycopg2

DB_USER = "qoqos"
DB_PASS = "StypticQtweDB"
DB_NAME = "qoqos_db"
DB_HOST = "localhost"
DB_PORT = "5432"

try:
    conn = psycopg2.connect(dbname=DB_NAME, user=DB_USER, password=DB_PASS, host=DB_HOST, port=DB_PORT)
    cur = conn.cursor()
    cur.execute("UPDATE \"BlogPost\" SET image = REPLACE(image, '/assets/', '/mock_pictures/');")
    conn.commit()
    print("Paths fixed successfully!")
except Exception as e:
    print(f"Error: {e}")
finally:
    if 'conn' in locals():
        conn.close()
