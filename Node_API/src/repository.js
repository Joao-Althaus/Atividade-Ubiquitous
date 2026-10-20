const pool = require('./database');

const columns = 'id, name, email, birthday, cpf, inserted_at, updated_at';

class UserRepository {
  async list() {
    const result = await pool.query(
      `SELECT ${columns} FROM users ORDER BY id`,
    );
    return result.rows;
  }

  async get(id) {
    const result = await pool.query(
      `SELECT ${columns} FROM users WHERE id = $1`,
      [id],
    );
    return result.rows[0] || null;
  }

  async create(user) {
    const query = `
      INSERT INTO users (name, email, birthday, cpf)
      VALUES ($1, $2, $3, $4)
      RETURNING ${columns}
    `;
    const values = [user.name, user.email, user.birthday, user.cpf];
    const result = await pool.query(query, values);
    return result.rows[0];
  }

  async update(id, user) {
    const query = `
      UPDATE users
      SET name = $1, email = $2, birthday = $3, cpf = $4,
          updated_at = CURRENT_TIMESTAMP
      WHERE id = $5
      RETURNING ${columns}
    `;
    const values = [user.name, user.email, user.birthday, user.cpf, id];
    const result = await pool.query(query, values);
    return result.rows[0] || null;
  }

  async delete(id) {
    const result = await pool.query('DELETE FROM users WHERE id = $1', [id]);
    return result.rowCount > 0;
  }
}

module.exports = UserRepository;
