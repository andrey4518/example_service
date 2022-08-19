FROM python:3.10.0-slim-buster
RUN apt update
RUN apt install -y git build-essential
RUN pip install -U pip
RUN pip install nipyapi
COPY ./app /app
COPY ./spec /spec
WORKDIR /app

CMD ["python", "main.py"]