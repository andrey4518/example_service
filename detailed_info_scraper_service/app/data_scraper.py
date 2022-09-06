from imdb import Cinemagoer
from tmdbv3api import TMDb
from tmdbv3api import Movie
import settings as s

class ImdbScraper:
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

    def __init__(self) -> None:
        self.ia = Cinemagoer()

    def collect_info(self, movie_data):
        movie_info = self.ia.get_movie(movie_data['imdb_id'])
        movie_info = {k: movie_info.get(k) for k in self.movie_data_keys}
        movie_info['movie_id'] = movie_data['id']
        for field in self.list_fields:
            if movie_info.get(field) is None:
                movie_info[field] = []
        for field in self.str_fields:
            if movie_info.get(field) is None:
                movie_info[field] = ''
        movie_info['original_title'] = movie_info['original title']
        del movie_info['original title']
        movie_info['plot_outline'] = movie_info['plot outline']
        del movie_info['plot outline']
        return movie_info


class TmdbScraper:
    def __init__(self) -> None:
        tmdb = TMDb()
        tmdb.api_key = s.get_tmdb_api_key()
        self.movie = Movie()

    def resolve_video_url(self, video_info):
        if video_info["site"] == "YouTube":
            return f'https://www.youtube.com/watch?v={video_info["key"]}'
        print(f'TMDB Scraper: Unsupported video site: {video_info["site"]}')
        return ""

    def collect_info(self, movie_data):
        m = self.movie.details(movie_data['tmdb_id'])
        return {
            "movie_id": movie_data['id'],
            "adult": m['adult'],
            "genres": [x['name'] for x in m['genres']],
            "homepage": m['homepage'],
            "original_title": m['original_title'],
            "overview": m['overview'],
            "popularity": m['popularity'],
            "runtime": m['runtime'],
            "tagline": m['tagline'],
            "title": m['title'],
            "vote_average": m['vote_average'],
            "vote_count": m['vote_count'],
            "keywords": 
                list(set(
                    sum(
                        [
                            list(map(lambda x: x.get('name', ''), keywords))
                            for _, keywords in m['keywords'].items()
                        ],
                        []
                    )
                )),
            "video_urls": [self.resolve_video_url(info) for info in m['videos']['results']],
        }
