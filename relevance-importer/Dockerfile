FROM python:3.9
WORKDIR app
ENV ELASTIC_URL "localhost"
ENV ELASTIC_INDEX "users"
ENV FILES_AND_RELEVANCE '/tmp/relevance_1.txt:99,/tmp/relevance_2.txt:50'
COPY . .
RUN pip install elasticsearch
CMD ["python", "./main.py"]