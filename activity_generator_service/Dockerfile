FROM python:3.10.0-slim-buster
RUN apt update
RUN apt install -y git build-essential
COPY requirements.txt /app/requirements.txt
RUN pip install -U pip
RUN pip install -r /app/requirements.txt
COPY ./app /app
WORKDIR /app

CMD ["python", "main.py"]
