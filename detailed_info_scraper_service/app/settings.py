import os

def get_movie_creation_topic_name():
    return os.getenv("MOVIE_CREATION_TOPIC_NAME", 'test-topic')

def get_kafka_url():
    return os.getenv("KAFKA_URL", 'kafka:9092')

def get_imdb_insert_api_path():
    return f'{os.getenv("API_URL", "http://api:8080/api/v1")}/movie_imdb_info'

def get_tmdb_insert_api_path():
    return f'{os.getenv("API_URL", "http://api:8080/api/v1")}/movie_tmdb_info'

def get_tmdb_api_key():
    return os.getenv("TMDB_API_V3_KEY")