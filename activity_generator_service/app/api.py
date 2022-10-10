from fastapi import APIRouter
import requests
from log import logger
from scheduler import get_scheduler
from dto import CurrentScheduledJobsResponse, JobCreateDeleteResponse
import json
import aiohttp
import local_db_queries as q
import os

router = APIRouter()

from faker import Faker
fake = Faker()

@router.get('/generate_user')
async def gen_user():
    profile = fake.profile([
            'username',
            'name',
            'sex',
            'address',
            'mail',
        ])
    profile['email'] = profile['mail']
    del profile['mail']
    return profile


async def generate_and_create_users(endpoint, count=1):
    users = [await gen_user() for _ in range(count)]
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(users)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{users}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))
            for u in (await r.json()).get('users', []):
                id = u.get('id')
                if id:
                    await q.set_user_service_id(
                        await q.get_random_user_id(),
                        id
                    )


@router.get('/start_user_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/users/insert_batch'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(generate_and_create_users,'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }

@router.get('/generate_users_batch')
async def generate_usesrs_batch(
    endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/users/insert_batch',
    count: int=100
    ):
    await generate_and_create_users(endpoint, count)
    return {'status': 'done'}


@router.get('/generate_movie')
async def gen_movie():
    movie = await q.get_random_movie()
    links = await q.get_links_by_movie_id(movie.id)
    genres = await q.get_genres_by_movie_id(movie.id)
    return {
        'id': movie.id,
        'name': movie.title,
        'imdb_id': links.imdb_id,
        'tmdb_id': links.tmdb_id,
        'genres': [g.name for g in genres]
    }


@router.get('/get_movies_ready_to_export')
async def get_movies_ready_to_export(count: int = 5):
    movies = await q.get_random_movies_ready_to_export(count)
    return [
        {
            'id': m.id,
            'name': m.title,
            'imdb_id': (await q.get_links_by_movie_id(m.id)).imdb_id,
            'tmdb_id': (await q.get_links_by_movie_id(m.id)).tmdb_id,
            'genres': [g.name for g in await q.get_genres_by_movie_id(m.id)]
        }
        for m in movies
    ]


async def generate_and_create_movies(endpoint, count=1):
    internal_movies = await get_movies_ready_to_export(count)
    internal_movie_ids = [m.pop('id') for m in internal_movies]
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(internal_movies)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{internal_movies}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))
            for internal_movie, id in zip(internal_movies, internal_movie_ids):
                internal_movie['id'] = id

            for external_movie in (await r.json()).get('movies', []):
                external_id = external_movie.get('id')
                if external_id:
                    internal_id = None
                    for internal_movie in internal_movies:
                        if internal_movie.get('name') == external_movie.get('name'):
                            internal_id = internal_movie.get('id')
                            break
                    if internal_id:
                        await q.set_movie_exported(
                            internal_id,
                            external_id
                        )


@router.get('/generate_movies_batch')
async def generate_usesrs_batch(
    endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/movies/insert_batch',
    count: int=100
    ):
    await generate_and_create_movies(endpoint, count)
    return {'status': 'done'}


@router.get('/start_movie_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/movies'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(generate_and_create_movies, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }


@router.get('/get_exported')
async def get_exported():
    return {
        'movies': await q.get_exported_movies(),
        'users': await q.get_exported_users(),
        'ratings': await q.get_exported_ratings(),
        'tags': await q.get_exported_tags()
    }


@router.get('/get_ratings_ready_to_export')
async def get_ratings_ready_to_export(count:int = 0):
    if count:
        return (await q.get_all_ratings_ready_to_export())[:count]
    return await q.get_all_ratings_ready_to_export()


async def generate_and_create_ratings(endpoint, count=1):
    internal_data = await get_ratings_ready_to_export(count)
    insert_ratings = [
        {
            'user_id': x['user'].service_id,
            'movie_id': x['movie'].service_id,
            'rating': x['rating'].rate,
        }
        for x in internal_data
    ]
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(insert_ratings)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{insert_ratings}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))

            for external_rating in (await r.json()).get('ratings', []):
                internal_id = None
                for internal_item in internal_data:
                    if all([
                        external_rating['user_id'] == internal_item['user'].service_id,
                        external_rating['movie_id'] == internal_item['movie'].service_id,
                    ]):
                        internal_id = internal_item['rating'].id
                        break
                if internal_id:
                    q.set_rating_exported(internal_id)


@router.get('/generate_ratings_batch')
async def generate_ratings_batch(
    endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/ratings/insert_batch',
    count: int=100
    ):
    await generate_and_create_ratings(endpoint, count)
    return {'status': 'done'}


@router.get('/start_rating_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/ratings'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(generate_and_create_ratings, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }


@router.get('/generate_tag_ready_to_export')
async def generate_tag_ready_to_export():
    return await q.get_tag_ready_to_export()


async def tag_generating_job(endpoint):
    tag = await generate_tag_ready_to_export()
    if not tag:
        return
    tag_id = tag.pop('id')
    tag.pop('exported')
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(tag)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{tag}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))
            json_data = await r.json()
            id = json_data.get('tag', {}).get('id')
            if id:
                await q.set_tag_exported(tag_id)


@router.get('/start_tag_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/tags'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(tag_generating_job, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }
