# gateway

O gerenciamento das APIs da solução é feito com o API Gateway [Kong](https://konghq.com/kong/), o que faz a centralização e exposição das APIs de busca e autenticação com base na configuração do arquivo `kong.yml`.


A configuração consisiste em:
- Exposição do serviço de busca de usuários no path `/search/v1`
- Exposição e transformação do serviço autenticação no path `/auth`, onde o usuário não precisa informar parâmetros internos como `client_id`, `client_secret` e `grant_type` para obter o token de autenticação
- Restrição do acesso no serviço de busco, onde é exigido token válido no header `Authorization: Bearer <token>` das requisições
- Configurações as níveis de APIs de rate limite e cache
