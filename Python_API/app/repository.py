from psycopg.rows import dict_row

from .database import POOL

COLUMNS = "id, name, email, birthday, cpf, inserted_at, updated_at"


class UserRepository:
    def list(self):
        with POOL.connection() as connection:
            with connection.cursor(row_factory=dict_row) as cursor:
                cursor.execute(f"SELECT {COLUMNS} FROM users ORDER BY id")
                return cursor.fetchall()

    def get(self, user_id: int):
        with POOL.connection() as connection:
            with connection.cursor(row_factory=dict_row) as cursor:
                cursor.execute(
                    f"SELECT {COLUMNS} FROM users WHERE id = %s",
                    (user_id,),
                )
                return cursor.fetchone()

    def create(self, data: dict):
        query = (
            f"INSERT INTO users (name, email, birthday, cpf) "
            f"VALUES (%(name)s, %(email)s, %(birthday)s, %(cpf)s) "
            f"RETURNING {COLUMNS}"
        )
        with POOL.connection() as connection:
            with connection.cursor(row_factory=dict_row) as cursor:
                cursor.execute(query, data)
                return cursor.fetchone()

    def update(self, user_id: int, data: dict):
        query = (
            f"UPDATE users "
            f"SET name = %(name)s, email = %(email)s, birthday = %(birthday)s, "
            f"cpf = %(cpf)s, updated_at = CURRENT_TIMESTAMP "
            f"WHERE id = %(id)s "
            f"RETURNING {COLUMNS}"
        )
        parameters = {**data, "id": user_id}
        with POOL.connection() as connection:
            with connection.cursor(row_factory=dict_row) as cursor:
                cursor.execute(query, parameters)
                return cursor.fetchone()

    def delete(self, user_id: int) -> bool:
        with POOL.connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute("DELETE FROM users WHERE id = %s", (user_id,))
                return cursor.rowcount > 0
