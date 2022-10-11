from log import logger
from scheduler import get_scheduler
from api.router import router
import aiohttp
import json
from dto import JobCreateDeleteResponse
from local_db import tag as tag_q
import os

@router.get('/get_tags_ready_to_export')
async def get_tags_ready_to_export(count:int = 0):
    if count:
        return (await tag_q.get_all_tags_ready_to_export())[:count]
    return await tag_q.get_all_tags_ready_to_export()


async def generate_and_create_tags(endpoint, count=1):
    internal_data = await get_tags_ready_to_export(count)
    insert_tags = [
        {
            'user_id': x['user'].service_id,
            'movie_id': x['movie'].service_id,
            'tag_text': x['tag'].tag_text,
        }
        for x in internal_data
    ]
    async with aiohttp.ClientSession() as session:
        async with session.post(endpoint, data=json.dumps(insert_tags)) as r:
            logger.info((
                f'Request to "{r.url}" with payload "{insert_tags}" finished '
                f'with code {r.status} and response "{await r.text()}"'
            ))

            for external_tag in (await r.json()).get('tags', []):
                internal_id = None
                for internal_item in internal_data:
                    if all([
                        external_tag['user_id'] == internal_item['user'].service_id,
                        external_tag['movie_id'] == internal_item['movie'].service_id,
                    ]):
                        internal_id = internal_item['tag'].id
                        break
                if internal_id:
                    tag_q.set_tag_exported(internal_id)


@router.get('/generate_tags_batch')
async def generate_tags_batch(
    endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/tags/insert_batch',
    count: int=100
    ):
    await generate_and_create_tags(endpoint, count)
    return {'status': 'done'}


@router.get('/start_tag_generating_activity',response_model=JobCreateDeleteResponse,tags=["generating"])
async def start_user_generating_activity(interval: int=5, endpoint: str=f'{os.getenv("API_URL", "http://api:8080/api/v1")}/tags/insert_batch'):
    scheduler = await get_scheduler()
    job = scheduler.add_job(generate_and_create_tags, 'interval', seconds=interval, args=[endpoint])
    return {
        "scheduled":True,
        "job_id":job.id
    }
