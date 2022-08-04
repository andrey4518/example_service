from sqlalchemy import Column, Integer, String, BigInteger, Boolean, Float, ForeignKey
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class User(Base):
    __tablename__ = 'user'

    id = Column(BigInteger().with_variant(Integer, "sqlite"), primary_key=True)
    original_id = Column(BigInteger, unique=True)
    service_id = Column(BigInteger)
    exported = Column(Boolean())

    rating = relationship("Rating", back_populates="user")
    tag = relationship("Tag", back_populates="user")


class Movie(Base):
    __tablename__ = 'movie'

    id = Column(BigInteger().with_variant(Integer, "sqlite"), primary_key=True)
    original_id = Column(BigInteger, unique=True)
    service_id = Column(BigInteger)
    title = Column(String())
    exported = Column(Boolean())

    movie_genre = relationship("MovieGenre", back_populates="movie")
    links = relationship("MovieLinks", back_populates="movie")
    rating = relationship("Rating", back_populates="movie")
    tag = relationship("Tag", back_populates="movie")


class Genre(Base):
    __tablename__ = 'genre'

    id = Column(BigInteger().with_variant(Integer, "sqlite"), primary_key=True)
    name = Column(String())
    exported = Column(Boolean())

    movie_genre = relationship("MovieGenre", back_populates="genre")


class MovieGenre(Base):
    __tablename__ = 'movie_genre'

    movie_id = Column(BigInteger().with_variant(Integer, "sqlite"), ForeignKey("movie.id"), primary_key=True)
    movie = relationship("Movie", back_populates="movie_genre")
    genre_id = Column(BigInteger().with_variant(Integer, "sqlite"), ForeignKey("genre.id"), primary_key=True)
    genre = relationship("Genre", back_populates="movie_genre")
    exported = Column(Boolean())


class MovieLinks(Base):
    __tablename__ = 'movie_links'

    id = Column(BigInteger().with_variant(Integer, "sqlite"), primary_key=True)
    movie_id = Column(BigInteger, ForeignKey("movie.id"))
    movie = relationship("Movie", back_populates="links")
    imdb_id = Column(BigInteger)
    tmdb_id = Column(BigInteger)
    exported = Column(Boolean())


class Rating(Base):
    __tablename__ = 'rating'

    id = Column(BigInteger().with_variant(Integer, "sqlite"), primary_key=True)
    user_id = Column(BigInteger, ForeignKey("user.id"))
    user = relationship("User", back_populates="rating")
    movie_id = Column(BigInteger, ForeignKey("movie.id"))
    movie = relationship("Movie", back_populates="rating")
    rate = Column(Float)
    exported = Column(Boolean())


class Tag(Base):
    __tablename__ = 'tag'

    id = Column(BigInteger().with_variant(Integer, "sqlite"), primary_key=True)
    user_id = Column(BigInteger, ForeignKey("user.id"))
    user = relationship("User", back_populates="tag")
    movie_id = Column(BigInteger, ForeignKey("movie.id"))
    movie = relationship("Movie", back_populates="tag")
    tag_text = Column(String)
    exported = Column(Boolean())
