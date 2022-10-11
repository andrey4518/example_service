from typing import List
import sqlalchemy as sa
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import func, select
from sqlalchemy.orm import sessionmaker
from local_db.common import create_local_db_engine
from local_db.dto import Movie, MovieGenre, MovieLinks, Rating, User, Genre

async def get_random_movies_ready_to_export(count=1):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Movie)\
            .join(Movie.rating)\
            .join(Rating.user)\
            .where(sa.and_(
                Movie.exported == False,
                User.exported
            )).order_by(func.random()).limit(count)
        result = await session.execute(stmt)
        movies = result.fetchall()
        if movies is None:
            return []
        else:
            movies = [m[0] for m in movies]
            movies = {m.id: m for m in movies}
            movies = list(movies.values())
            movies = movies[:count]
        return movies


async def get_links_by_movie_id(movie_id) -> MovieLinks:
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(MovieLinks).where(MovieLinks.movie_id == movie_id)
        result = await session.execute(stmt)
        links = result.fetchone()[0]

    return links


async def get_genres_by_movie_id(movie_id) -> List[Genre]:
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Genre.name).join(Genre.movie_genre)\
            .where(MovieGenre.movie_id == movie_id)
        result = await session.execute(stmt)
        genres = [row for row in result]

    return genres


async def set_movie_exported(movie_id, service_id):
    engine = await create_local_db_engine()
    async_session = sessionmaker(
        engine, expire_on_commit=False, class_=AsyncSession
    )

    async with async_session() as session:
        stmt = select(Movie, MovieLinks, MovieGenre)\
            .join(Movie.links)\
            .join(Movie.movie_genre)\
            .where(Movie.id == movie_id)
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
        stmt = select(Movie).where(Movie.exported)
        result = await session.execute(stmt)
        return [row.Movie for row in result]
