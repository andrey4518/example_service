from aiokafka import AIOKafkaConsumer
import asyncio
import json
from imdb import Cinemagoer
import requests
import os


async def consume():
    ia = Cinemagoer()
    movie_data_keys = [
        'genres', # str
        'original title', # str
        'runtimes', # list(str)
        'countries', # list(str)
        'rating', # float
        'votes', # int
        'plot outline', # str
        'languages', # list(str)
        'year', # int
        'kind', # str
        'plot', # list(str)
        'synopsis', # list(str)
    ]
    list_fields = [
        'runtimes',
        'countries',
        'languages',
        'plot',
        'synopsis',
    ]
    str_fields = [
        'genres',
        'original title',
        'plot outline',
        'kind',
    ]
    print('creating consumer')
    consumer = AIOKafkaConsumer(
        os.getenv("MOVIE_CREATION_TOPIC_NAME", 'test-topic'),
        bootstrap_servers=os.getenv("KAFKA_URL", 'kafka:9092')
    )
    # Get cluster layout and join group `my-group`
    print('starting consumer')
    await consumer.start()
    try:
        # Consume messages
        print('wait for messages')
        async for msg in consumer:
            print("consumed: ", msg.topic, msg.partition, msg.offset,
                  msg.key, msg.value, msg.timestamp)
            data = json.loads(msg.value)
            print(f'data: {data}')
            if data.get('type') == 'Movie':
                movie_info = ia.get_movie(data['value']['imdb_id'])
                movie_info = {k: movie_info.get(k) for k in movie_data_keys}
                movie_info['movie_id'] = data['value']['id']
                for field in list_fields:
                    if movie_info.get(field) is None:
                        movie_info[field] = []
                for field in str_fields:
                    if movie_info.get(field) is None:
                        movie_info[field] = ''
                movie_info['original_title'] = movie_info['original title']
                del movie_info['original title']
                movie_info['plot_outline'] = movie_info['plot outline']
                del movie_info['plot outline']
                r = requests.post(
                    f'{os.getenv("API_URL", "http://api:8080/api/v1")}/movie_imdb_info',
                    data=json.dumps(movie_info)
                )
                print(movie_info)
                if r.status_code != 200:
                    print(f'got error {r.text}')
    except Exception as e:
        print('got exception: ', e)
    finally:
        # Will leave consumer group; perform autocommit if enabled.
        print('stoping consumer')
        await consumer.stop()


if __name__ == '__main__':
    print('starting...')
    asyncio.run(consume())
    print('finishing...')