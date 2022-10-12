from sqlalchemy.ext.asyncio import create_async_engine
from async_lru import alru_cache


@alru_cache
async def create_local_db_engine():
    return create_async_engine(
        "sqlite+aiosqlite:///data/local.db",
        future=True
    )
