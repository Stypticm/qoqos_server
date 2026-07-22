const { Client } = require('pg');
const fs = require('fs');

const client = new Client({
  user: 'qoqos',
  host: 'localhost',
  database: 'qoqos_db',
  password: 'StypticQtweDB',
  port: 5432,
});

async function exportData() {
  try {
    await client.connect();
    console.log('Connected to Postgres');
    await client.query("SET client_encoding TO 'UTF8'");

    // 1. Блог
    const blogRes = await client.query('SELECT * FROM "BlogPost"');
    const blog_posts = blogRes.rows;
    console.log('First blog title:', blog_posts[0]?.title);

    // 2. Лоты
    const lotRes = await client.query('SELECT * FROM "MarketplaceLot"');
    const lots = lotRes.rows;

    // 3. Заявки (Skupka)
    const skupkaRes = await client.query('SELECT * FROM "Skupka"');
    const skupkas = skupkaRes.rows;

    // 4. Заказы (Order)
    const orderRes = await client.query('SELECT * FROM "Order"');
    const orders = orderRes.rows;

    // 5. Мастера (Master)
    const masterRes = await client.query('SELECT * FROM "Master"');
    const masters = masterRes.rows;

    // 6. Персонал (User)
    const staffRes = await client.query('SELECT * FROM "User"');
    const staff = staffRes.rows;

    // 7. Пользователи (User)
    const usersRes = await client.query('SELECT * FROM "User"');
    const users = usersRes.rows;

    const data = {
      blog_posts,
      lots,
      skupkas,
      orders,
      masters,
      staff,
      users
    };

    fs.writeFileSync('db_data.json', JSON.stringify(data, null, 2));
    console.log('Successfully exported data to db_data.json');

  } catch (err) {
    console.error('Error exporting data:', err);
  } finally {
    await client.end();
  }
}

exportData();
