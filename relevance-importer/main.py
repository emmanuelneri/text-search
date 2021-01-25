import os

from elasticsearch import Elasticsearch

field = 'relevance'

files_and_relevance = os.getenv('FILES_AND_RELEVANCE', 'relevance_1.txt:99,relevance_2.txt:50')
elasticsearch_host = os.getenv('ELASTIC_HOST', 'localhost')
index = os.getenv('ELASTIC_INDEX', 'users')
es = Elasticsearch([{'host': elasticsearch_host, 'port': 9200}])

print(f'starting relevance importer with configs: elasticsearch: {elasticsearch_host} - index: {index}')


def update_relevance(filename, relevance):
    print(f"starting update relevance started file: {filename} - relevance: {relevance}")
    with open(filename, 'r') as f:
        for line in f:
            doc = {
                'query': {
                    'match_phrase': {
                        'ID': line.strip()
                    },
                },
                '_source': 'false'
            }
            res = es.search(index=index, body=doc)
            doc_id = res['hits']['hits'][0]['_id']
            es.update(index=index, id=doc_id, body={"doc": {field: relevance}})


for file_and_relevance in files_and_relevance.split(','):
    split = file_and_relevance.split(':')
    file = split[0]
    relevance = split[1]
    update_relevance(file, relevance)

es.indices.refresh(index=index)

