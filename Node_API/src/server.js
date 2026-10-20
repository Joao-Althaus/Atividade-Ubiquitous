const express = require('express');
const pool = require('./database');
const { router, errorHandler } = require('./controller');

const app = express();
app.use(express.json({ limit: '32kb' }));
app.get('/health', (_request, response) => response.json({ status: 'ok' }));
app.use(router);
app.use((_request, response) => response.status(404).json({ error: 'not found' }));
app.use(errorHandler);

const port = Number(process.env.PORT || 8080);

async function waitForDatabase() {
  for (let attempt = 0; attempt < 30; attempt += 1) {
    try {
      await pool.query('SELECT 1');
      return;
    } catch (error) {
      if (attempt === 29) {
        throw error;
      }
      console.log('waiting for postgres');
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }
  }
}

async function startServer() {
  await waitForDatabase();
  app.listen(port, '0.0.0.0', () => {
    console.log(`Node API listening on ${port}`);
  });
}

startServer().catch((error) => {
  console.error(error);
  process.exit(1);
});
