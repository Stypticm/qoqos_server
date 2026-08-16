require('dotenv').config();
const { Telegraf, Markup } = require('telegraf');
const bcrypt = require('bcryptjs');
const { Client } = require('pg');
const crypto = require('crypto');

const bot = new Telegraf(process.env.TELEGRAM_BOT_TOKEN);

// Конфигурация БД
const dbConfig = {
  user: process.env.DB_USER || 'qoqos',
  host: 'localhost',
  database: process.env.DB_NAME || 'qoqos_db',
  password: process.env.DB_PASSWORD,
  port: 5432,
};

const MASTER_PASSWORD = process.env.MASTER_PASSWORD;

// Состояния пользователей (в памяти для простоты, можно в Redis)
const userStates = new Map();

async function getDbClient() {
  const client = new Client(dbConfig);
  await client.connect();
  return client;
}

// Генерация случайного пароля
function generatePassword(length = 10) {
  return crypto.randomBytes(length).toString('base64').slice(0, length).replace(/\//g, 'x').replace(/\+/g, 'y');
}

bot.start((ctx) => {
  ctx.reply('🤖 Бот для управления учетными записями Qoqos\n\nОтправьте Telegram ID (для сотрудников) или Логин (для клиентов) чтобы управлять аккаунтом.');
});

bot.on('text', async (ctx) => {
  const text = ctx.message.text;
  const chatId = ctx.chat.id;
  const state = userStates.get(chatId);

  // 0. Обработка Reply-сообщений для чата на сайте
  if (ctx.message.reply_to_message) {
    const replyTo = ctx.message.reply_to_message;
    const originalText = (replyTo.text || replyTo.caption || '').trim();

    if (originalText.includes('Новый вопрос в чате от пользователя')) {
      const match = originalText.match(/от пользователя (\S+)/);
      if (match && match[1]) {
        const guestId = match[1];
        console.log('[ChatReply] Reply detected for guestId:', guestId, 'chatId:', chatId, 'replyToMessageId:', replyTo.message_id);
        const client = await getDbClient();
        try {
          await client.query('CREATE TABLE IF NOT EXISTS "ChatReply" ( id TEXT PRIMARY KEY, "guestId" TEXT, text TEXT, "isRead" BOOLEAN DEFAULT FALSE, "createdAt" TIMESTAMP DEFAULT NOW() )');
          const id = 'reply_' + Date.now();
          await client.query('INSERT INTO "ChatReply" (id, "guestId", text) VALUES ($1, $2, $3)', [id, guestId, text]);
          try {
            await ctx.reply(`✅ Ответ моментально отправлен пользователю ${guestId}!`);
          } catch (replyError) {
            console.error('[ChatReply] Failed to reply to admin:', replyError);
          }
        } catch (e) {
          console.error('[ChatReply] Error:', e);
          try {
            await ctx.reply('❌ Ошибка при отправке ответа.');
          } catch (replyError) {
            console.error('[ChatReply] Failed to reply error to admin:', replyError);
          }
        } finally {
          await client.end();
        }
        return;
      }
    }
    console.log('[ChatReply] Ignoring reply - wrong target message. Reply text:', originalText.substring(0, 100));
  }

  // 1. Проверка мастер-пароля
  if (state?.waitingForMaster) {
    if (text === MASTER_PASSWORD) {
      ctx.reply('✅ Доступ разрешен! С возвращением.\n\nИспользуйте /start для меню.');
      userStates.set(chatId, { ...state, authenticated: true, waitingForMaster: false });
      
      // После ввода мастер-пароля создаем аккаунт
      await createAccount(ctx, state.pendingUser);
    } else {
      ctx.reply('❌ Неверный пароль доступа.');
    }
    return;
  }

  // 2. Обработка ввода ID/Логина
  if (!state || !state.authenticated) {
    const client = await getDbClient();
    try {
      const res = await client.query('SELECT * FROM "User" WHERE "telegramId" = $1', [text]);
      if (res.rows.length > 0) {
        const user = res.rows[0];
        const newPass = generatePassword();
        const salt = bcrypt.genSaltSync(10);
        const hash = bcrypt.hashSync(newPass, salt);
        
        await client.query('UPDATE "User" SET "passwordHash" = $1 WHERE id = $2', [hash, user.id]);
        
        ctx.reply(`✅ Пароль успешно изменен!\n\n📱 Логин: ${user.telegramId}\n🔑 Новый пароль: ${newPass}\n\n⚠️ Сохраните новый пароль!`);
      } else {
        // Пользователь не найден - предлагаем выбрать роль
        userStates.set(chatId, { pendingId: text });
        ctx.reply('👤 Пользователь не найден. Выберите роль:', Markup.inlineKeyboard([
          [Markup.button.callback('👑 Admin', 'role_ADMIN'), Markup.button.callback('🔧 Master', 'role_MASTER')],
          [Markup.button.callback('📊 Manager', 'role_MANAGER'), Markup.button.callback('👤 User', 'role_USER')]
        ]));
      }
    } finally {
      await client.end();
    }
    return;
  }
});

bot.action(/role_(.+)/, async (ctx) => {
  const role = ctx.match[1];
  const chatId = ctx.chat.id;
  const state = userStates.get(chatId);

  if (!state || !state.pendingId) return;

  userStates.set(chatId, { ...state, pendingRole: role, waitingForMaster: true });
  ctx.reply('🔒 Бот защищен. Введите пароль доступа:');
});

async function createAccount(ctx, pendingUser) {
  const chatId = ctx.chat.id;
  const state = userStates.get(chatId);
  const telegramId = state.pendingId;
  const role = state.pendingRole;

  const client = await getDbClient();
  try {
    const password = generatePassword();
    const salt = bcrypt.genSaltSync(10);
    const hash = bcrypt.hashSync(password, salt);
    const id = 'user_' + Date.now();

    await client.query(
      'INSERT INTO "User" (id, "telegramId", "passwordHash", role, "createdAt", "updatedAt") VALUES ($1, $2, $3, $4, NOW(), NOW())',
      [id, telegramId, hash, role]
    );

    ctx.reply(`✅ Аккаунт успешно создан!\n\n📱 Логин/ID: ${telegramId}\n🔑 Пароль: ${password}\n👤 Роль: ${role}\n\n⚠️ Сохраните эти данные! Пароль больше не будет показан.`);
  } catch (e) {
    console.error(e);
    ctx.reply('❌ Ошибка при создании аккаунта.');
  } finally {
    await client.end();
    userStates.delete(chatId);
  }
}

bot.launch().then(() => console.log('🚀 Telegram Bot started'));

// Enable graceful stop
process.once('SIGINT', () => bot.stop('SIGINT'));
process.once('SIGTERM', () => bot.stop('SIGTERM'));

// -------------------------------------------------------------------
// Лёгкий HTTP-сервер для отдачи ответов в Next.js без установки 'pg'
// -------------------------------------------------------------------
const http = require('http');

const server = http.createServer(async (req, res) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Content-Type', 'application/json');

  if (req.method === 'GET' && req.url.startsWith('/poll?guestId=')) {
    const url = new URL(req.url, 'http://localhost');
    const guestId = url.searchParams.get('guestId');
    
    if (!guestId) {
      res.writeHead(400);
      return res.end(JSON.stringify({ error: 'guestId required' }));
    }

    const client = await getDbClient();
    try {
      await client.query('CREATE TABLE IF NOT EXISTS "ChatReply" ( id TEXT PRIMARY KEY, "guestId" TEXT, text TEXT, "isRead" BOOLEAN DEFAULT FALSE, "createdAt" TIMESTAMP DEFAULT NOW() )');
      
      const result = await client.query('SELECT id, text, "createdAt" FROM "ChatReply" WHERE "guestId" = $1 AND "isRead" = false ORDER BY "createdAt" ASC', [guestId]);
      const messages = result.rows;

      if (messages.length > 0) {
        const ids = messages.map(m => m.id);
        await client.query('UPDATE "ChatReply" SET "isRead" = true WHERE id = ANY($1::text[])', [ids]);
      }

      res.writeHead(200);
      res.end(JSON.stringify({ messages }));
    } catch (e) {
      console.error('[HTTP] DB Error:', e);
      res.writeHead(500);
      res.end(JSON.stringify({ error: 'DB Error' }));
    } finally {
      await client.end();
    }
  } else {
    res.writeHead(404);
    res.end(JSON.stringify({ error: 'Not found' }));
  }
});

server.on('error', (err) => {
  if (err.code === 'EADDRINUSE') {
    console.log('[HTTP] Port 8002 already in use - skipping HTTP server (bot may already be running)');
  } else {
    console.error('[HTTP] Server error:', err);
  }
});

server.listen(8002, () => console.log('[HTTP] Local API started on port 8002'));
