from fastapi import APIRouter
import requests
from log import logger
from scheduler import get_scheduler
from dto import CurrentScheduledJobsResponse, JobCreateDeleteResponse
import json
import aiohttp
import local_db_queries as q

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


async def user_generating_job(endpoint):
    user = await gen_user()
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(user)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{user}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))
            id = (await r.json()).get('id')
            if id:
                await q.set_user_service_id(
                    await q.get_random_user_id(),
                    id
                )


@router.get('/start_user_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str='http://example_service_api_1:8080/api/v1/users'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(user_generating_job,'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }


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


async def movie_generating_job(endpoint):
    movie = await gen_movie()
    movie_id = movie.pop('id')
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(movie)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{movie}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))
            json_data = await r.json()
            id = json_data.get('movie', {}).get('id')
            if id:
                await q.set_movie_exported(
                    movie_id,
                    id
                )


@router.get('/start_movie_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str='http://example_service_api_1:8080/api/v1/movies'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(movie_generating_job, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }


@router.get('/get_exported')
async def get_exported():
    return {
        'movies': await q.get_exported_movies(),
        'users': await q.get_exported_users(),
        'ratings': await q.get_exported_ratings()
    }


@router.get('/generate_rating_ready_to_export')
async def generate_rating_ready_to_export():
    return await q.get_rating_ready_to_export()


@router.get('/get_ratings_ready_to_export')
async def get_ratings_ready_to_export():
    return await q.get_all_ratings_ready_to_export()


async def rating_generating_job(endpoint):
    rating = await generate_rating_ready_to_export()
    if not rating:
        return
    rating_id = rating.pop('id')
    rating.pop('exported')
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(rating)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{rating}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))
            json_data = await r.json()
            id = json_data.get('rating', {}).get('id')
            if id:
                await q.set_rating_exported(rating_id)


@router.get('/start_rating_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str='http://example_service_api_1:8080/api/v1/ratings'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(rating_generating_job, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }
