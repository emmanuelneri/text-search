# user-importer

Projeto com objetivo de importar usuários através de arquivos csv para o Elasticsearch.

###### Modelo esperado no csv
```
6e172695-c76c-4364-8dd9-44e6d2d3aed9,Heitor Rovaron,heitor.rovaron
4e8660b0-7350-4211-9b9b-9ba50792ccd9,Melony Terci,melony.terci
```
Observação: O arquivo é esperado sem header

###### Formato enviado para o elasticsearch
```
{
   "ID":"6e172695-c76c-4364-8dd9-44e6d2d3aed9",
   "Name":"Heitor Rovaron",
   "Username":"heitor.rovaron",
   "relevance":0
}
{
   "ID":"6e172695-c76c-4364-8dd9-44e6d2d3aed9",
   "Name":"Melony Terci",
   "Username":"melony.terci",
   "relevance":0
}
```
Observação: Inicialmente todo usuário é inserido no index com relevância 0, onde posteriormente o campo será atualizado de acordo com os arquivos processados pela aplicação `relevance-importer`

#### Dependências
- `pip3 install elasticsearch`
- `pip3 install pandas`
