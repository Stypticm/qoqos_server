SET client_encoding = 'UTF8';

DELETE FROM "BlogPost";
DELETE FROM "MarketplaceLot";

INSERT INTO "BlogPost" (id, title, content, excerpt, image, category, author, published, "createdAt", "updatedAt") VALUES
('b1b2b3b4-1111-1111-1111-111111111111', 'Как подготовить iPhone к продаже', 'Перед тем как передать свой iPhone новому владельцу...', 'Полная инструкция по подготовке iPhone к выкупу', 'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5?q=80&w=600&auto=format&fit=crop', 'Инструкции', 'QOQOS Team', true, NOW(), NOW()),
('b1b2b3b4-2222-2222-2222-222222222222', 'Тренды рынка б/у электроники в 2026 году', 'Рынок вторичной электроники продолжает расти...', 'Анализ рынка б/у техники: почему восстановленные смартфоны становятся популярнее', 'https://images.unsplash.com/photo-1531403009284-440f080d1e12?q=80&w=600&auto=format&fit=crop', 'Аналитика', 'Илья Воронов', true, NOW(), NOW()),
('b1b2b3b4-3333-3333-3333-333333333333', 'AI и диагностика: как агенты оценивают устройства', 'В QOQOS мы используем систему ИИ-агентов...', 'Как технологии искусственного интеллекта помогают автоматизировать оценку', 'https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?q=80&w=600&auto=format&fit=crop', 'Технологии', 'AI Engineer', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO "MarketplaceLot" (id, title, model, storage, color, condition, description, price, photos, "coverPhoto", status, "telegramId", "sellerName", "createdAt", "updatedAt") VALUES
('m1m2m3m4-1111-1111-1111-111111111111', 'iPhone 14 Pro 128GB Space Black', 'iPhone 14 Pro', '128GB', 'Space Black', 'Отличное', 'Полностью протестирован нашими AI-агентами. Батарея 89%.', 72000.0, '{"https://images.unsplash.com/photo-1695048133142-1a20484d2569?q=80&w=500&auto=format&fit=crop"}', 'https://images.unsplash.com/photo-1695048133142-1a20484d2569?q=80&w=500&auto=format&fit=crop', 'available', 'system', 'Магазин QOQOS', NOW(), NOW()),
('m1m2m3m4-2222-2222-2222-222222222222', 'MacBook Air 13" M2 8/256GB Midnight', 'MacBook Air M2', '256GB', 'Midnight', 'Как новый', 'Всего 42 цикла заряда батареи. Родная клавиатура.', 95000.0, '{"https://images.unsplash.com/photo-1517336714731-489689fd1ca8?q=80&w=500&auto=format&fit=crop"}', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?q=80&w=500&auto=format&fit=crop', 'available', 'system', 'Магазин QOQOS', NOW(), NOW()),
('m1m2m3m4-3333-3333-3333-333333333333', 'AirPods Pro 2nd Gen (Type-C)', 'AirPods Pro 2', 'N/A', 'Белый', 'Хорошее', 'Звук чистый, шумоподавление работает идеально.', 16500.0, '{"https://images.unsplash.com/photo-1588449668365-d15e397f6787?q=80&w=500&auto=format&fit=crop"}', 'https://images.unsplash.com/photo-1588449668365-d15e397f6787?q=80&w=500&auto=format&fit=crop', 'available', 'system', 'Магазин QOQOS', NOW(), NOW())
ON CONFLICT DO NOTHING;
