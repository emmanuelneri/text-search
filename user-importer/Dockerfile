FROM python:3.9
WORKDIR app
ENV ELASTIC_URL "localhost"
ENV ELASTIC_INDEX "users"
ENV FILE_CHUNK_SIZE 10000
ENV FILE_NAME 'users.csv.gz'
COPY . .
RUN pip install pandas elasticsearch
CMD ["python", "./main.py"]