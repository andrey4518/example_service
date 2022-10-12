from log import logger
from scheduler import get_scheduler
from api.router import router
import aiohttp
import json
from dto import JobCreateDeleteResponse
from local_db import rating as rating_q
import os

@router.get('/get_ratings_ready_to_export')
async def get_ratings_ready_to_export(count:int = 0):
    if count:
        return (await rating_q.get_all_ratings_ready_to_export())[:count]
    return await rating_q.get_all_ratings_ready_to_export()


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
                    rating_q.set_rating_exported(internal_id)


@router.get('/generate_ratings_batch')
async def generate_ratings_batch(
    endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/ratings/insert_batch',
    count: int=100
    ):
    await generate_and_create_ratings(endpoint, count)
    return {'status': 'done'}


@router.get('/start_rating_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_rating_generating_activity(interval: int=5, endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/ratings/insert_batch'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(generate_and_create_ratings, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }
