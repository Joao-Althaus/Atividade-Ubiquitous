from contextlib import asynccontextmanager

from fastapi import FastAPI

from .controller import router
from .database import POOL


@asynccontextmanager
async def lifespan(_app):
    POOL.open(wait=True)
    yield
    POOL.close()


app = FastAPI(title="Users API - Python", version="1.0.0", lifespan=lifespan)
app.include_router(router)


@app.get("/health")
def health():
    return {"status": "ok"}
