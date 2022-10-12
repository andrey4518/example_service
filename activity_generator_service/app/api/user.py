from log import logger
from scheduler import get_scheduler
from api.router import router
import aiohttp
import json
from dto import JobCreateDeleteResponse
from local_db import user as user_q
import os

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
                    await user_q.set_user_service_id(
                        await user_q.get_random_user_id(),
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
