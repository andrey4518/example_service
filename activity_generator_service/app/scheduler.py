from time import sleep
from fastapi import APIRouter
from dto import CurrentScheduledJobsResponse, JobCreateDeleteResponse
from functools import lru_cache
from async_lru import alru_cache
from log import logger
#APScheduler Related Libraries
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.jobstores.sqlalchemy import SQLAlchemyJobStore


router = APIRouter()

@alru_cache()
async def get_scheduler() -> AsyncIOScheduler:
    try:
        jobstores = {
            'default': SQLAlchemyJobStore(url='sqlite:///jobs.sqlite')
        }
        Schedule = AsyncIOScheduler(jobstores=jobstores)
        Schedule.start()
        logger.info("Created Schedule Object")
        return Schedule
    except:
        logger.error("Unable to Create Schedule Object")
        raise


@router.get("/schedule/show_schedules/",response_model=CurrentScheduledJobsResponse,tags=["schedule"])
async def get_scheduled_syncs():
    """
    Will provide a list of currently Scheduled Tasks
    """
    schedules = []
    scheduler = await get_scheduler()
    for job in scheduler.get_jobs():
        schedules.append({"job_id": str(job.id), "run_frequency": str(job.trigger), "next_run": str(job.next_run_time)})
    return {"jobs":schedules}


@router.delete("/schedule",response_model=JobCreateDeleteResponse,tags=["schedule"])
async def remove_job_scheduler(job_id):
    """
    Remove a Job from a Schedule
    """
    try:
        scheduler = await get_scheduler()
        scheduler.remove_job(job_id)
    except Exception:
        pass
    return {"scheduled":False,"job_id":job_id}
