import psycopg2
import uuid
from datetime import datetime

# Данные подключения к БД (из твоего .env)
DB_USER = "qoqos"
DB_PASS = "StypticQtweDB"
DB_NAME = "qoqos_db"
DB_HOST = "localhost"
DB_PORT = "5432"

def seed_data():
    try:
        conn = psycopg2.connect(
            dbname=DB_NAME,
            user=DB_USER,
            password=DB_PASS,
            host=DB_HOST,
            port=DB_PORT
        )
        # Устанавливаем кодировку клиента в UTF-8, чтобы избежать кракозябр
        conn.set_client_encoding('UTF8')
        cur = conn.cursor()

        print("🧹 Очистка старых некорректных записей блога и товаров...")
        # Удаляем записи, где в названии много вопросительных знаков (проблемы с кодировкой)
        cur.execute('DELETE FROM "BlogPost" WHERE title LIKE \'%???%\';')
        cur.execute('DELETE FROM "MarketplaceLot" WHERE title LIKE \'%???%\';')

        # Пример данных для блога
        blog_posts = [
            (
                str(uuid.uuid4()),
                "Как подготовить iPhone к продаже: пошаговое руководство от QOQOS",
                "Перед тем как передать свой iPhone новому владельцу или сдать в программу Buyback, крайне важно правильно подготовить устройство. В этой статье мы расскажем, как сделать резервную копию в iCloud, отключить функцию 'Локатор', выйти из Apple ID и безопасно стереть все персональные данные.",
                "Полная инструкция по подготовке iPhone к выкупу: бэкап, отвязка от iCloud и сброс настроек.",
                "https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?q=80&w=600&auto=format&fit=crop",
                "Инструкции",
                "QOQOS Team",
                True,
                datetime.utcnow(),
                datetime.utcnow()
            ),
            (
                str(uuid.uuid4()),
                "Тренды рынка б/у электроники в 2026 году: что покупают чаще всего?",
                "Рынок вторичной электроники продолжает расти рекордными темпами. Экологическая осознанность и экономическая выгода мотивируют людей выбирать восстановленные устройства. В топ продаж этого сезона входят iPhone 14 Pro, MacBook Air на чипах M2 и AirPods Pro 2. Рассказываем, почему спрос смещается в сторону восстановленных девайсов с гарантией.",
                "Анализ рынка б/у техники: почему восстановленные смартфоны и ноутбуки становятся популярнее новых.",
                "https://images.unsplash.com/photo-1531403009284-440f080d1e12?q=80&w=600&auto=format&fit=crop",
                "Аналитика",
                "Илья Воронов",
                True,
                datetime.utcnow(),
                datetime.utcnow()
            ),
            (
                str(uuid.uuid4()),
                "Искусственный интеллект и диагностика: как AI-агенты оценивают устройства",
                "Как сделать оценку б/у смартфона максимально честной и быстрой? В QOQOS мы используем систему ИИ-агентов. pricing-agent парсит актуальные рыночные цены, а diagnostics-agent анализирует состояние внутренних датчиков смартфона через специальный тест-клиент. Это позволяет рассчитать справедливую цену за 3 минуты.",
                "Как технологии искусственного интеллекта помогают автоматизировать оценку электроники.",
                "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?q=80&w=600&auto=format&fit=crop",
                "Технологии",
                "AI Engineer",
                True,
                datetime.utcnow(),
                datetime.utcnow()
            )
        ]

        print("📝 Вставка тестовых записей в таблицу BlogPost...")
        for post in blog_posts:
            cur.execute(
                '''
                INSERT INTO "BlogPost" (id, title, content, excerpt, image, category, author, published, "createdAt", "updatedAt")
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (id) DO NOTHING;
                ''',
                post
            )

        # Пример данных для Marketplace (Lots)
        lots = [
            (
                str(uuid.uuid4()),
                "iPhone 14 Pro 128GB Space Black (Отличное состояние)",
                "iPhone 14 Pro",
                "128GB",
                "Космический черный",
                "Apple",
                "SKU-IP14P-128-SB",
                "Отличное",
                "Полностью протестирован нашими AI-агентами. Батарея 89%. На экране нет царапин, корпус без сколов. В комплекте идет оригинальная коробка и новый кабель Type-C - Lightning. Предоставляется гарантия 90 дней от сервиса QOQOS.",
                72000.0,
                79000.0,
                ["https://images.unsplash.com/photo-1695048133142-1a20484d2569?q=80&w=500&auto=format&fit=crop"],
                "https://images.unsplash.com/photo-1695048133142-1a20484d2569?q=80&w=500&auto=format&fit=crop",
                "available",
                "system",
                "Магазин QOQOS",
                12,
                False,
                datetime.utcnow(),
                datetime.utcnow()
            ),
            (
                str(uuid.uuid4()),
                "MacBook Air 13\" M2 8/256GB Midnight (Как новый)",
                "MacBook Air M2",
                "256GB",
                "Темная ночь (Midnight)",
                "Apple",
                "SKU-MBA-M2-256-MN",
                "Как новый",
                "Состояние нового ноутбука. Всего 42 цикла заряда батареи (ёмкость 98%). Родная клавиатура с русской гравировкой. Без следов использования. Полный комплект с зарядным устройством 30W Dual Port.",
                95000.0,
                105000.0,
                ["https://images.unsplash.com/photo-1517336714731-489689fd1ca8?q=80&w=500&auto=format&fit=crop"],
                "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?q=80&w=500&auto=format&fit=crop",
                "available",
                "system",
                "Магазин QOQOS",
                28,
                False,
                datetime.utcnow(),
                datetime.utcnow()
            ),
            (
                str(uuid.uuid4()),
                "AirPods Pro 2nd Gen (Type-C) (Хорошее состояние)",
                "AirPods Pro 2",
                "N/A",
                "Белый",
                "Apple",
                "SKU-APP2-TC",
                "Хорошее",
                "Наушники полностью дезинфицированы и очищены. Кейс имеет мелкие царапины от ношения без чехла. Звук чистый, активное шумоподавление и режим прозрачности работают идеально. Зарядный порт USB Type-C.",
                16500.0,
                19900.0,
                ["https://images.unsplash.com/photo-1588449668365-d15e397f6787?q=80&w=500&auto=format&fit=crop"],
                "https://images.unsplash.com/photo-1588449668365-d15e397f6787?q=80&w=500&auto=format&fit=crop",
                "available",
                "system",
                "Магазин QOQOS",
                5,
                False,
                datetime.utcnow(),
                datetime.utcnow()
            )
        ]

        print("📝 Вставка тестовых записей в таблицу MarketplaceLot...")
        for lot in lots:
            cur.execute(
                '''
                INSERT INTO "MarketplaceLot" (
                    id, title, model, storage, color, brand, sku, condition, description, 
                    price, "oldPrice", photos, "coverPhoto", status, "telegramId", "sellerName", 
                    "viewsCount", "isAccessory", "createdAt", "updatedAt"
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (id) DO NOTHING;
                ''',
                lot
            )

        conn.commit()
        print("✅ База данных успешно заполнена качественными русскоязычными mock-данными!")

    except Exception as e:
        print(f"❌ Ошибка заполнения БД: {e}")
    finally:
        if 'conn' in locals():
            conn.close()

if __name__ == "__main__":
    seed_data()
