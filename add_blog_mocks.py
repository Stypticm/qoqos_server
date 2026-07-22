import psycopg2
import uuid
from datetime import datetime

DB_USER = "qoqos"
DB_PASS = "StypticQtweDB"
DB_NAME = "qoqos_db"
DB_HOST = "localhost"
DB_PORT = "5432"

try:
    conn = psycopg2.connect(dbname=DB_NAME, user=DB_USER, password=DB_PASS, host=DB_HOST, port=DB_PORT)
    conn.set_client_encoding('UTF8')
    cur = conn.cursor()

    blog_posts = [
        (
            str(uuid.uuid4()),
            "Дорси запускает Buzz — конкурента Slack с открытым исходным кодом и встроенными ИИ-агентами",
            "В во вторник Джек Дорси представил Buzz — бесплатную платформу для совместной работы с открытым исходным кодом, разработанную компанией Block, которая объединяет людей и ИИ-агентов в общих рабочих пространствах. Инструмент уже доступен на buzz.xyz для macOS, Windows и Linux и позиционируется как децентрализованная альтернатива Slack и GitHub.",
            "Джек Дорси представил Buzz — новую децентрализованную альтернативу Slack с ИИ-агентами.",
            "/assets/pic_1.jpg",
            "Технологии",
            "Новости IT",
            True,
            datetime.utcnow(),
            datetime.utcnow()
        ),
        (
            str(uuid.uuid4()),
            "Wildberries предлагает кредиты, а не выплаты, продавцам, пострадавшим от ударов украинских дронов",
            "В субботу, 18 июля, украинские дроны атаковали два крупных логистических центра Wildberries — крупнейшего российского онлайн-маркетплейса. В результате ударов погибли восемь человек, более 70 получили ранения. Это один из самых разрушительных ударов по российской коммерческой инфраструктуре с начала войны.",
            "Wildberries предлагает кредиты продавцам после атак на логистические центры.",
            "/assets/pic_2.jpg",
            "Новости",
            "Новости",
            True,
            datetime.utcnow(),
            datetime.utcnow()
        ),
        (
            str(uuid.uuid4()),
            "Китайский WAIC представил масштабные модели ИИ, встревожив Уолл-стрит",
            "Всемирная конференция по искусственному интеллекту 2026 года завершилась в Шанхае в воскресенье, подведя итог четырём дням, наглядно продемонстрировавшим стремительное взросление китайской отрасли ИИ. Кульминацией стали презентации масштабных открытых моделей и программная речь председателя Си Цзиньпина, представившего Пекин как архитектора нового глобального миропорядка в сфере ИИ. Череда громких анонсов взбудоражила рынки: в пятницу индекс Nasdaq Composite упал на 1,4% на фоне распродажи акций полупроводниковых компаний.",
            "Всемирная конференция по искусственному интеллекту 2026 в Шанхае.",
            "/assets/pic_3.jpg",
            "Искусственный Интеллект",
            "Тех-аналитика",
            True,
            datetime.utcnow(),
            datetime.utcnow()
        )
    ]

    for post in blog_posts:
        cur.execute(
            '''
            INSERT INTO "BlogPost" (id, title, content, excerpt, image, category, author, published, "createdAt", "updatedAt")
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (id) DO NOTHING;
            ''',
            post
        )

    conn.commit()
    print("Mock data for blog successfully inserted!")
except Exception as e:
    print(f"Error: {e}")
finally:
    if 'conn' in locals():
        conn.close()
