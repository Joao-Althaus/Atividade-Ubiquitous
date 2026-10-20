import os

from psycopg_pool import ConnectionPool

HOST = os.getenv("DB_HOST", "localhost")
PORT = os.getenv("DB_PORT", "5432")
NAME = os.getenv("DB_NAME", "users_db")
USER = os.getenv("DB_USER", "users")
PASSWORD = os.getenv("DB_PASSWORD", "users_dev_password")

POOL = ConnectionPool(
    f"host={HOST} port={PORT} dbname={NAME} user={USER} password={PASSWORD}",
    min_size=1,
    max_size=10,
    open=False,
)
