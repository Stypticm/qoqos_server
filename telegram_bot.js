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

// Мастер-пароль из скриншота
const MASTER_PASSWORD = 'GolyanovoRomaMisha';

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
