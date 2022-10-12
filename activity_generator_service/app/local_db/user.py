from typing import Optional
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import func, select
from sqlalchemy.orm import sessionmaker
from local_db.common import create_local_db_engine
from local_db.dto import User


async def get_random_user_id():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(User.id).where(User.exported == False)\
            .order_by(func.random()).limit(1)
        result = await session.execute(stmt)
        id = result.fetchone()
        if id is None:
            id = None
        else:
            id = id[0]
    return id


async def set_user_service_id(user_id: int, service_id: Optional[int]):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(User).where(User.id == user_id)
        result = await session.execute(stmt)
        user: User = result.scalars().first()
        user.service_id = service_id
        user.exported = service_id is not None
        await session.commit()


async def get_exported_users():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(User).where(User.exported)
        result = await session.execute(stmt)
        return [row.User for row in result]
