from sqlalchemy import create_engine, select
from sqlalchemy.orm import Session
from local_db.dto import Base, Genre, Movie, MovieGenre, Rating, Tag, User, MovieLinks
import pandas as pd
from itertools import chain
from csv import DictReader
from tqdm import tqdm

engine = create_engine(
    "sqlite:///data/local.db",
    # echo=True,
    future=True
)

Base.metadata.create_all(engine)

def sqlalchemy_custom_batched_add(objects, session: Session, batch_size: int):
    for idx, obj in enumerate(objects):
        session.add(obj)
        if idx % batch_size == 0:
            session.commit()
    session.commit()

with Session(engine) as session:
    print('Reading movies.csv')
    df = pd.read_csv('../../data/ml-25m/movies.csv')
    print('Prepare genres')
    df['genres'] = df['genres'].apply(lambda x: x.split('|'))
    genres = set(chain.from_iterable(df['genres'].values))
    genres = list(tqdm(map(lambda x: Genre(name=x, exported=False), genres)))

    print('Prepare movies')
    genres_mapping = {g.name: g for g in genres}

    movies = [
        (
            Movie(original_id=row['movieId'], service_id=None, title=row['title'], exported=False),
            row['genres']
        )
        for _, row in df.iterrows()
    ]

    for movie, genres_list in tqdm(movies):
        for genre_name in genres_list:
            genre = genres_mapping[genre_name]
            mg = MovieGenre(exported=False)
            mg.genre = genre
            movie.movie_genre.append(mg)
            session.add(mg)

    print('Commit movies and genres')
    for m, _ in movies:
        session.add(m)

    session.add_all(genres)

    session.commit()

    movies_original_id_mapping = {m.original_id: m for m, _ in movies}

    print('Clear some memory')
    del movies
    del genres
    del genres_mapping
    del df

    print('Fill links')
    with open('../../data/ml-25m/links.csv') as csvfile:
        reader = DictReader(csvfile)
        links = [
            MovieLinks(
                movie=movies_original_id_mapping[int(row['movieId'])],
                imdb_id=row['imdbId'],
                tmdb_id=row['tmdbId'],
                exported=False
                )
            for row in tqdm(reader)
        ]
        print('Commit links')
        session.add_all(links)
        session.commit()
        print('Clear a little bit')
        del links

    print('Fill users')
    with open('../../data/ml-25m/ratings.csv') as ratings_csvfile:
        with open('../../data/ml-25m/tags.csv') as tags_csvfile:
            ratings_reader = DictReader(ratings_csvfile)
            tags_reader = DictReader(tags_csvfile)
            user_ids = set([
                int(row['userId'])
                for row in tqdm(chain.from_iterable([ratings_reader, tags_reader]))
            ])
            users = [User(original_id=id, service_id=None, exported=False) for id in tqdm(user_ids)]
            
            print('Commit users')
            session.add_all(users)

            session.commit()

            print('Clear a little bit')
            users_original_id_mapping = {u.original_id: u for u in users}

            del user_ids
            del users

    print('Inserting tags')

    with open('../../data/ml-25m/tags.csv') as csvfile:
        sqlalchemy_custom_batched_add(
            map(
                lambda row: Tag(
                    tag_text=row['tag'],
                    exported=False,
                    movie=movies_original_id_mapping[int(row['movieId'])],
                    user=users_original_id_mapping[int(row['userId'])],
                ),
                tqdm(DictReader(csvfile))
            ),
            session,
            batch_size=500_000
        )

    print('Inserting ratings')

    with open('../../data/ml-25m/ratings.csv') as ratings_csvfile:
        sqlalchemy_custom_batched_add(
            map(
                lambda row: Rating(
                    rate=float(row['rating']),
                    exported=False,
                    movie=movies_original_id_mapping[int(row['movieId'])],
                    user=users_original_id_mapping[int(row['userId'])],
                ),
                tqdm(DictReader(ratings_csvfile))
            ),
            session,
            batch_size=500_000
        )
