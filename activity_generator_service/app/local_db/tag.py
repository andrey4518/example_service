from typing import Optional
import sqlalchemy as sa
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import func, select
from sqlalchemy.orm import sessionmaker
from local_db.common import create_local_db_engine
from local_db.dto import Tag, Movie, User


async def get_all_tags_ready_to_export():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Tag, Movie, User)\
            .join(Tag.movie)\
            .join(Tag.user)\
            .where(sa.and_(
                Tag.exported == False,
                Movie.exported,
                User.exported
            ))
        result = await session.execute(stmt)
        return [
            {'tag': row[0], 'movie': row[1], 'user': row[2]}
            for row in result
        ]


async def get_exported_tags():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Tag).where(Tag.exported)
        result = await session.execute(stmt)
        return [row.Tag for row in result]


async def set_tag_exported(tag_id):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Tag).where(Tag.id == tag_id)
        result = await session.execute(stmt)
        tag = result.fetchone()[0]
        tag.exported = True
        await session.commit()
