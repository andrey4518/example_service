import asyncio
from typing import List, Optional
import sqlalchemy as sa
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy import func, select
from sqlalchemy.orm import sessionmaker
import local_db_dto
from async_lru import alru_cache
from log import logger


@alru_cache
async def create_local_db_engine():
    return create_async_engine(
        "sqlite+aiosqlite:///data/local.db",
        future=True
    )


async def get_random_user_id():
    engine = create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.User.id).where(local_db_dto.User.exported == False)\
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
        stmt = select(local_db_dto.User).where(local_db_dto.User.id == user_id)
        result = await session.execute(stmt)
        user: local_db_dto.User = result.scalars().first()
        user.service_id = service_id
        user.exported = service_id is not None
        await session.commit()


async def get_exported_users():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.User).where(local_db_dto.User.exported)
        result = await session.execute(stmt)
        return [row.User for row in result]


async def get_random_movie() -> local_db_dto.Movie:
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Movie).where(local_db_dto.Movie.exported == False)\
            .order_by(func.random()).limit(1)
        result = await session.execute(stmt)
        movie = result.fetchone()[0]

    return movie


async def get_links_by_movie_id(movie_id) -> local_db_dto.MovieLinks:
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.MovieLinks).where(local_db_dto.MovieLinks.movie_id == movie_id)
        result = await session.execute(stmt)
        links = result.fetchone()[0]

    return links


async def get_genres_by_movie_id(movie_id) -> List[local_db_dto.Genre]:
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Genre.name).join(local_db_dto.Genre.movie_genre)\
            .where(local_db_dto.MovieGenre.movie_id == movie_id)
        result = await session.execute(stmt)
        genres = [row for row in result]

    return genres


async def set_movie_exported(movie_id, service_id):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Movie, local_db_dto.MovieLinks, local_db_dto.MovieGenre)\
            .join(local_db_dto.Movie.links)\
            .join(local_db_dto.Movie.movie_genre)\
            .where(local_db_dto.Movie.id == movie_id)
        for r in await session.execute(stmt):
            r.Movie.service_id = service_id
            r.Movie.exported = service_id is not None
            r.MovieLinks.exported = service_id is not None
            r.MovieGenre.exported = service_id is not None
        await session.commit()


async def get_exported_movies():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Movie).where(local_db_dto.Movie.exported)
        result = await session.execute(stmt)
        return [row.Movie for row in result]


async def get_rating_ready_to_export():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Rating)\
            .join(local_db_dto.Rating.movie)\
            .join(local_db_dto.Rating.user)\
            .where(sa.and_(
                local_db_dto.Rating.exported == False,
                local_db_dto.Movie.exported,
                local_db_dto.User.exported
            )).order_by(func.random()).limit(1)
        result = await session.execute(stmt)
        rating = result.fetchone()
        if rating is None:
            return None
        else:
            rating = rating[0]
        return rating


async def get_all_ratings_ready_to_export():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Rating)\
            .join(local_db_dto.Rating.movie)\
            .join(local_db_dto.Rating.user)\
            .where(sa.and_(
                local_db_dto.Rating.exported == False,
                local_db_dto.Movie.exported,
                local_db_dto.User.exported
            ))
        result = await session.execute(stmt)
        return [row.Rating for row in result]


async def get_exported_ratings():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Rating).where(local_db_dto.Rating.exported)
        result = await session.execute(stmt)
        return [row.Rating for row in result]


async def set_rating_exported(rating_id):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Rating).where(local_db_dto.Rating.id == rating_id)
        result = await session.execute(stmt)
        rating = result.fetchone()[0]
        rating.exported = True
        await session.commit()


async def get_tag_ready_to_export():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Tag)\
            .join(local_db_dto.Tag.movie)\
            .join(local_db_dto.Tag.user)\
            .where(sa.and_(
                local_db_dto.Tag.exported == False,
                local_db_dto.Movie.exported,
                local_db_dto.User.exported
            )).order_by(func.random()).limit(1)
        result = await session.execute(stmt)
        tag = result.fetchone()
        if tag is None:
            return None
        else:
            tag = tag[0]
        return tag


async def get_exported_tags():
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Tag).where(local_db_dto.Tag.exported)
        result = await session.execute(stmt)
        return [row.Tag for row in result]


async def set_tag_exported(tag_id):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(local_db_dto.Tag).where(local_db_dto.Tag.id == tag_id)
        result = await session.execute(stmt)
        tag = result.fetchone()[0]
        tag.exported = True
        await session.commit()
