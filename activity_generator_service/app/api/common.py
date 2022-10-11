from api.router import router
from local_db import movie, user, rating, tag

@router.get('/get_exported')
async def get_exported():
    return {
        'movies': await movie.get_exported_movies(),
        'users': await user.get_exported_users(),
        'ratings': await rating.get_exported_ratings(),
        'tags': await tag.get_exported_tags()
    }