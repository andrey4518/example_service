import sqlalchemy as sa
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from sqlalchemy.orm import sessionmaker
from local_db.common import create_local_db_engine
from local_db.dto import Rating, Movie, User

async def get_all_ratings_ready_to_export():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Rating, Movie, User)\
            .join(Rating.movie)\
            .join(Rating.user)\
            .where(sa.and_(
                Rating.exported == False,
                Movie.exported,
                User.exported
            ))
        result = await session.execute(stmt)
        return [
            {'rating': row[0], 'movie': row[1], 'user': row[2]}
            for row in result
        ]


async def get_exported_ratings():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Rating).where(Rating.exported)
        result = await session.execute(stmt)
        return [row.Rating for row in result]


async def set_rating_exported(rating_id):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Rating).where(Rating.id == rating_id)
        result = await session.execute(stmt)
        rating = result.fetchone()[0]
        rating.exported = True
        await session.commit()
