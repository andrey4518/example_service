from fastapi import FastAPI
import uvicorn
from scheduler import router as SchedulerRouter, get_scheduler
from api.router import router as ApiRouter
from log import logger


app = FastAPI(
    title="Test data generation service",
    description="App for generating test data (registration of new user, adding content, user behavior and so on)",
    version="0.0.1",
)

app.include_router(SchedulerRouter)
app.include_router(ApiRouter)

@app.on_event("startup")
async def load_schedule_or_create_blank():
    """
    Instatialise the Schedule Object as a Global Param and also load existing Schedules from SQLite
    This allows for persistent schedules across server restarts. 
    """
    get_scheduler()


@app.on_event("shutdown")
async def pickle_schedule():
    """
    An Attempt at Shutting down the schedule to avoid orphan jobs
    """
    get_scheduler().shutdown()
    logger.info("Disabled Schedule")

if __name__ == '__main__':
    uvicorn.run(app, host='0.0.0.0', port=8000, debug=True)
