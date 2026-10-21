# Users APIs: Go, Python, Node.js e Java

Projeto acadêmico com quatro APIs REST de CRUD para usuários, implementadas em Go, Python, Node.js e Java com Spring Boot. Cada API possui Dockerfile próprio e as quatro compartilham PostgreSQL no Docker Compose.

## Recursos de usuários

As quatro APIs oferecem os mesmos endpoints. A porta muda conforme a implementação:

| API | URL base |
| --- | --- |
| Go | `http://localhost:8081` |
| Python | `http://localhost:8082` |
| Node.js | `http://localhost:8083` |
| Java | `http://localhost:8084` |

```http
GET     /users       → Lista todos os usuários
GET     /users/{id}  → Consulta um usuário
POST    /users       → Cria um usuário
PUT     /users/{id}  → Atualiza um usuário
DELETE  /users/{id}  → Exclui um usuário
GET     /health      → Verifica se a API está disponível
```

O corpo de `POST` e `PUT` recebe os mesmos campos:

```json
{
  "name": "Ana Souza",
  "email": "ana@example.com",
  "birthday": "2000-01-31",
  "cpf": "12345678901"
}
```

As respostas de usuário incluem `id`, `name`, `email`, `birthday`, `cpf`, `inserted_at` e `updated_at`. O identificador e as datas são gerados pelo banco. Email e CPF são únicos; o CPF é armazenado como texto para preservar zeros à esquerda.

| Operação | Resposta de sucesso | Erros principais |
| --- | --- | --- |
| `GET /users` | `200` e uma lista | `500` |
| `GET /users/{id}` | `200` e um usuário | `400` para ID inválido, `404` se não encontrado |
| `POST /users` | `201` e o usuário criado | `400` para dados inválidos, `409` para email/CPF duplicado |
| `PUT /users/{id}` | `200` e o usuário atualizado | `400`, `404` ou `409` |
| `DELETE /users/{id}` | `204` sem corpo | `400` ou `404` |
| `GET /health` | `200` com `{"status":"ok"}` | — |

Os erros são retornados no formato `{"error":"mensagem"}`. O nome deve ter de 1 a 120 caracteres, o email até 254 caracteres, a data usa `YYYY-MM-DD` e o CPF contém 11 dígitos.

## Tecnologias

- **Go** com `pgx`
- **Python** com FastAPI e Psycopg
- **Node.js** com Express e `pg`
- **Java** com Spring Boot e JDBC
- **PostgreSQL** como banco de dados compartilhado
- **pgAdmin** para administração visual do banco
- **Docker Compose** para construir e iniciar os serviços

## Estrutura do projeto

```text
Atividade-Ubiquitous/
├── GO_API/
├── Python_API/
├── Node_API/
├── Java_API/
├── database/
├── postman/Users-APIs.postman_collection.json
├── docker-compose.yml
└── README.md
```

As APIs separam o recebimento de requisições, as regras de negócio e o acesso a dados em controller, service e repository. No Go, esses componentes podem ficar organizados no pacote `internal/user`.
## Requisitos

- Docker Desktop instalado e iniciado usando o engine Linux.
- Docker Compose disponível (`docker compose version`).
- Portas locais 5432, 5050 e 8081–8084 livres.

Não é necessário instalar Go, Python, Node.js, Java ou Maven no host para subir as imagens.

## Subir os serviços

Abra o PowerShell na pasta que contém `docker-compose.yml` e valide a configuração:

```powershell
docker info
docker compose config
```

Construa as imagens e inicie PostgreSQL, pgAdmin e as quatro APIs:

```powershell
docker compose up --build
```

A primeira execução pode levar alguns minutos para baixar imagens e dependências. Mantenha esse terminal aberto para acompanhar os logs. Para rodar em segundo plano:

```powershell
docker compose up --build -d
docker compose ps
```

Para acompanhar logs em segundo plano:

```powershell
docker compose logs -f
docker compose logs -f go-api python-api node-api java-api postgres
```

URLs locais:

| Serviço | Endereço |
| --- | --- |
| API Go | http://localhost:8081 |
| API Python | http://localhost:8082 |
| API Node.js | http://localhost:8083 |
| API Java | http://localhost:8084 |
| pgAdmin | http://localhost:5050 |
| PostgreSQL | localhost:5432 |

Verifique se cada API está pronta:

```powershell
8081,8082,8083,8084 | ForEach-Object { Invoke-RestMethod "http://localhost:$($_)/health" }
```

## Postman

A coleção pronta está em `postman/Users-APIs.postman_collection.json`.

1. No Postman, selecione **Import**.
2. Escolha o arquivo JSON acima.
3. Abra uma pasta de API (Go, Python, Node.js ou Java).
4. Execute **Health** e **Listar usuários** para conferir a API.
5. Para o CRUD, execute na ordem **Criar usuário**, **Buscar usuário por ID**, **Atualizar usuário** e **Excluir usuário**. A chamada de criação salva automaticamente o ID retornado na variável da coleção daquela API.
6. A chamada **Consultar usuário inexistente (404)** demonstra o caso de erro.

As quatro APIs usam o mesmo PostgreSQL. A coleção usa email e CPF distintos por API para que seja possível executar o fluxo em todas elas sem conflito de unicidade.

## Acessar o banco pelo navegador (pgAdmin)

Abra http://localhost:5050 e faça login com:

- Email: `admin@example.com`
- Senha: `admin`

O servidor **Users PostgreSQL** é provisionado pelo arquivo `database/servers.json`. Se for necessário registrar uma conexão manualmente, use estes parâmetros:

| Campo | Valor |
| --- | --- |
| Host | `postgres` |
| Porta | `5432` |
| Maintenance database | `users_db` |
| Usuário | `users` |
| Senha | `users_dev_password` |
| SSL | Prefer (ou desligado para ambiente local) |

Dentro do pgAdmin, o host é `postgres` porque ele acessa o banco pela rede interna do Compose. `localhost` é o endereço usado por ferramentas instaladas diretamente no Windows.

Para consultar os dados: **Servers → Users PostgreSQL → Databases → users_db → Schemas → public → Tables → users**. Clique com o botão direito na tabela e escolha **Query Tool**, ou execute:

```sql
SELECT id, name, email, birthday, cpf, inserted_at, updated_at
FROM users
ORDER BY id;
```

## Parar e reiniciar

Parar os serviços e manter os dados do banco:

```powershell
docker compose down
```

Iniciar novamente:

```powershell
docker compose up -d
```

Reconstruir as imagens sem usar cache:

```powershell
docker compose build --no-cache
docker compose up -d
```

Para remover também os dados persistidos no volume PostgreSQL (operação destrutiva):

```powershell
docker compose down -v
```

O schema em `database/init.sql` é aplicado automaticamente somente quando o volume do banco é criado pela primeira vez. Apagar o volume remove os dados e permite inicializar novamente o schema.
