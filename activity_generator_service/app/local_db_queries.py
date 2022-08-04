import asyncio
from typing import List, Optional
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy import func, select
from sqlalchemy.orm import sessionmaker, aliased
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
    engine = create_local_db_engine()
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
        logger.info(stmt)
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
