# relevance-importer

Projeto com objetivo de atualizar a relevância dos usuários através de arquivos txt para o Elasticsearch.

###### Modelo esperado no text
```
fba0be35-7111-43c5-8111-b326360da4d0
7354ff5e-cc72-4cc7-a8d0-279f3349c52b
4096545a-3d93-476d-9a25-ae486a12a720
```
Observação: O arquivo é esperado sem header

###### Envio para o elasticsearch
O Envio para elasticsearch é feito item a item atualizando a relavância de acordo com o valor passado junto ao arquivo na configuração `FILES_AND_RELEVANCE`.

Por exemplo, caso informado a configuração abaixo, os IDs do arquivo `relevance_1.txt` terão relevância 99 e os IDs do arquivo `relevance_2.txt` terão relevância 50.
```
-e FILES_AND_RELEVANCE="/tmp/relevance_1.txt:99,/tmp/relevance_2.txt:50"
```

O valor da relavância tem influência no score do registro nas buscas aplicadas no Elasticsearch, onde o valor desse campo é utilizado pela função `field_value_factor` da busca. 

#### Dependências
- `pip3 install elasticsearch`
