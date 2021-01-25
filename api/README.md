# api

Projeto com objetivo de disponibilizar as APIs de consulta de usuário.


APIs:
- GET /users?search=Diether
```
{
   "users":[
      {
         "ID":"fba0be35-7111-43c5-8111-b326360da4d0",
         "Name":"Diether Bein",
         "Username":"diether.bein"
      }
   ],
   "scrollId":"FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFnNDRmhtLVZ1U0N5U1ZPZm9mRXVicFEAAAAAAAACHxZ4UHlTamNkUVJzQ1hITkVfZjlzei1R"
}
```
As pesquisas são limitadas em até 15 registros e a ordenação do resultado é aplicado com base no campo `relevance` do registro, onde quanto maior o valor do campo relevância, mais bem classificado o registro será na busca. 

- GET /users/{id}/scroll
```
{
   "users":[
      {
         "ID":"1a2093a9-0bdb-43c7-acf3-f2d2decb9f25",
         "Name":"Jamerson Heidemann Aparecidi",
         "Username":"jamersonheidemannaparecidi"
      }
   ],
   "scrollId":"FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFnNDRmhtLVZ1U0N5U1ZPZm9mRXVicFEAAAAAAAACHhZ4UHlTamNkUVJzQ1hITkVfZjlzei1R"
}
```
A paginação por scroll não necessita dos parâmetros de busca para realizar a navegação nos registros, com isso, apenas é necessário informar o `scrollId` que é retornado durantes as consultas. 

- GET /health
```
{
    "status":"UP"
}
```
Retorna o status da aplicação, onde é realizado apenas uma chamda de função para validar se a aplicação consegue atender uma requisição http.

- GET /health/ready
```
{
   "status":"UP"
}
```
Retorna se o status da aplicação está pronta para uso, com isso, é realizado uma teste de conexão com o Elasticsearch para validar a saúde das dependências da aplicação.


###### Testes

Execute o comando para rodar os testes unitários
```
go test -v ./...
```

Execute o comando para rodar os testes de integração
```
go test -v -run ^TestIntegration$ ./...
```

###### Documentação
Execute o comando para gerar a documentação das APIs na pasta `doc`.
```
swag init -g internal/http/router.go -o doc
```

