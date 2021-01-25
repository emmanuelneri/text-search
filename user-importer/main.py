import gzip
import os

import pandas as pd
from elasticsearch import Elasticsearch, helpers

elasticsearch_host = os.getenv('ELASTIC_HOST', 'localhost')
index = os.getenv('ELASTIC_INDEX', 'users')
chunk_size = os.getenv('FILE_CHUNK_SIZE', 1000)
file_name = os.getenv('FILE_NAME', 'users.csv.gz')

es = Elasticsearch([{'host': elasticsearch_host, 'port': 9200}], timeout=30, max_retries=3, retry_on_timeout=True)

print(f'starting user importer with configs: elasticsearch: {elasticsearch_host} - index: {index} '
      f'- file: {file_name} - chunkSize: {chunk_size}')

with gzip.open(file_name, 'rb') as f:
    csv_file = pd.read_csv(f, delimiter=',', iterator=True, chunksize=int(chunk_size),
                           header=3, names=['ID', 'Name', 'Username'])
    for i, df in enumerate(csv_file):
        df["relevance"] = 0
        records = df.where(pd.notnull(df), None).T.to_dict()
        list_records = [records[it] for it in records]
        helpers.bulk(client=es, actions=list_records, index=index)
        print(f'records sent. chunk: {i}')
    es.indices.refresh(index=index)
