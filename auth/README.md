# auth

A autenticação da solução é realizada pelo [Keycloak](https://www.keycloak.org/), o qual fornece um serviço "out of the box" para realizar autenticações no padrão OAuth2 configurados no arquivo `keycloak-realm-api.json`.

A configuração consisiste em:
- Criação de um realm específico para o API Gateway
- Criação de um usuário padrão para acessar a API de busca de usuários

credenciais do usuário:
```
username: api@gmail.com.br
password: api
```
- Criação deu client específico para a aplicação de busca de usuários (`api`) e configurações do client para possibilitar autenticação por `grant_type=password`. 

credenciais do client:
```
"clientId": "api",
"secret": "1677c713-9509-45c5-8dcd-f466e842a9fa",
```
