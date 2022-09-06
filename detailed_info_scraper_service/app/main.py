from aiokafka import AIOKafkaConsumer
import asyncio
import json
from imdb import Cinemagoer
import requests
import settings as s
from data_scraper import ImdbScraper, TmdbScraper

async def consume():
    imdb_scraper = ImdbScraper()
    tmdb_scraper = TmdbScraper()
    print('creating consumer')
    consumer = AIOKafkaConsumer(
        s.get_movie_creation_topic_name(),
        bootstrap_servers=s.get_kafka_url()
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
                try:
                    imdb_info = imdb_scraper.collect_info(data['value'])
                    r = requests.post(
                        s.get_imdb_insert_api_path(),
                        data=json.dumps(imdb_info)
                    )
                    print(imdb_info)
                    if r.status_code != 200:
                        print(f'got error {r.text}')
                except Exception as e:
                    print('got exception: ', e)

                try:
                    tmdb_info = tmdb_scraper.collect_info(data['value'])
                    r = requests.post(
                        s.get_tmdb_insert_api_path(),
                        data=json.dumps(tmdb_info)
                    )
                    print(tmdb_info)
                    if r.status_code != 200:
                        print(f'got error {r.text}')
                except Exception as e:
                    print('got exception: ', e)
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