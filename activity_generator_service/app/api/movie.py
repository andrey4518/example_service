from log import logger
from scheduler import get_scheduler
from api.router import router
import aiohttp
import json
from dto import JobCreateDeleteResponse
from local_db import movie as movie_q
import os

@router.get('/get_movies_ready_to_export')
async def get_movies_ready_to_export(count: int = 5):
    movies = await movie_q.get_random_movies_ready_to_export(count)
    return [
        {
            'id': m.id,
            'name': m.title,
            'imdb_id': (await movie_q.get_links_by_movie_id(m.id)).imdb_id,
            'tmdb_id': (await movie_q.get_links_by_movie_id(m.id)).tmdb_id,
            'genres': [g.name for g in await movie_q.get_genres_by_movie_id(m.id)]
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
                        await movie_q.set_movie_exported(
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
async def start_user_generating_activity(interval: int=5, endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/movies/insert_batch'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(generate_and_create_movies, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }
