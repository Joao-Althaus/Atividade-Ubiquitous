# Users APIs: Go, Python, Node.js e Java

Projeto acadêmico com quatro APIs REST que implementam o mesmo CRUD de usuários. As aplicações compartilham um PostgreSQL e iniciam juntas pelo Docker Compose. Cada API possui seu próprio Dockerfile e uma porta externa.

## Recursos de usuários

As quatro APIs oferecem os mesmos endpoints:

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

As respostas de usuário incluem `id`, `name`, `email`, `birthday`, `cpf`, `inserted_at` e `updated_at`. O identificador e as datas de criação e atualização são gerados pelo banco de dados. Email e CPF são únicos; o CPF é armazenado como texto para manter zeros à esquerda.

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
Tarefa APIs/
├── GO_API/
│   ├── internal/user/
│   │   ├── controller.go
│   │   ├── model.go
│   │   ├── postgres_repository.go
│   │   ├── service.go
│   │   └── validation.go
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── Python_API/
│   ├── app/
│   │   ├── controller.py
│   │   ├── database.py
│   │   ├── main.py
│   │   ├── repository.py
│   │   ├── schemas.py
│   │   └── service.py
│   ├── Dockerfile
│   └── requirements.txt
├── Node_API/
│   ├── src/
│   │   ├── controller.js
│   │   ├── database.js
│   │   ├── repository.js
│   │   ├── server.js
│   │   └── service.js
│   ├── Dockerfile
│   └── package.json
├── Java_API/
│   ├── src/main/java/edu/users/
│   │   ├── controller/UserController.java
│   │   ├── model/
│   │   ├── repository/UserRepository.java
│   │   ├── service/UserService.java
│   │   └── UsersApiApplication.java
│   ├── Dockerfile
│   └── pom.xml
├── database/
│   ├── Dockerfile
│   ├── init.sql
│   └── servers.json
├── docker-compose.yml
└── README.md
```

As APIs separam o recebimento de requisições, as regras de negócio e o acesso a dados em controller, service e repository. No Go, esses componentes ficam no pacote `internal/user`.

## Configuração do banco de dados

O `docker-compose.yml` configura automaticamente o PostgreSQL e as APIs. Os valores padrão para desenvolvimento local são:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=users_db
DB_USER=users
DB_PASSWORD=users_dev_password
DB_SSLMODE=disable
```

Dentro da rede do Compose, as APIs usam `postgres` como host do banco. O schema é criado pelo arquivo `database/init.sql` na primeira inicialização do volume de dados.

## Iniciar a aplicação

Com Docker Desktop aberto e o daemon ativo, execute na raiz do projeto:

```sh
docker compose up --build
```

Endereços locais:

| Serviço | Endereço |
| --- | --- |
| Go API | http://localhost:8081 |
| Python API | http://localhost:8082 |
| Node.js API | http://localhost:8083 |
| Java API | http://localhost:8084 |
| pgAdmin | http://localhost:5050 |
| PostgreSQL | `localhost:5432` |

Acesso ao pgAdmin: `admin@example.com` / `admin`. O servidor PostgreSQL já aparece configurado no pgAdmin. As credenciais do banco são `users` / `users_dev_password`, com o banco `users_db`.

O volume `postgres_data` mantém os dados entre reinicializações. Para remover o volume e recriar o banco do zero, use `docker compose down -v`.

## Exemplo de uso

O exemplo abaixo cria um usuário e lista os usuários na API Go. Para testar outra implementação, substitua a porta por `8082`, `8083` ou `8084`.

```sh
curl -X POST http://localhost:8081/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ana Souza","email":"ana@example.com","birthday":"2000-01-31","cpf":"12345678901"}'

curl http://localhost:8081/users
```

No PowerShell do Windows, use `curl.exe` para chamar o executável curl em vez do alias `curl`.

## Execução local sem Docker

É necessário iniciar e configurar o PostgreSQL separadamente e definir as variáveis de conexão no ambiente. Versões usadas pelas imagens do Compose:

- Go 1.23
- Python 3.12
- Node.js 22
- Java 21 e Maven

A execução por Docker Compose é o caminho recomendado para iniciar o banco e as quatro APIs com a mesma configuração.
