# Technical requirements
Build a microservice to handle companies. It should provide the following operations:
- Create
- Patch
- Delete
- Get (one)

Only authenticated users should have access to create, update and delete companies.

# Domain

Each company is defined by the following attributes:

- ID (uuid) required
- Name (15 characters) required - unique
- Description (3000 characters) optional
- Amount of Employees (int) required
- Registered (boolean) required
- Type (Corporations | NonProfit | Cooperative | Sole Proprietorship) required


# Nice to have

Will be considered a plus:
- The solution to be production ready.
- On each mutating operation, an event should be produced.
- Dockerize the application to be ready for building the production docker image
- Use docker for setting up the external services such as the database
- REST is suggested, but GRPC is also an option
- JWT for authentication
- Kafka for events
- DB is up to you
- Integration tests are highly appreciated
- Linter
- Configuration file

# Prerequisite

1. Clone project

```sh
git clone https://github.com/stuton/xm-golang-exercise.git
```

2. Up and running docker containers

```sh
make run
```

3. Execute database migrations

```sh
make migrate
```

4. Create kafka cluster in kafka-ui (localhost:9093), where:

- Cluster name: xm-golang
- Bootstrap Servers: **kafka** with port **9092**

## Examples API calls

Get JWT token for following API calls using default admin credentials

```sh
curl -X POST localhost:8080/api/v1/login -H 'Content-Type: application/json' -d '{"username":"admin","password":"admin"}'
```

Create company, please don't forget to set JWT token

```sh
curl -X POST localhost:8080/api/v1/companies \
   -H 'Content-Type: application/json' \
   -H 'Authorization: JWT_TOKEN' \
   -d '{"name":"company1","amount_employees":10,"registered":false,"type":"Corporations"}'
```

You can find message in Kafka using kafka-ui (localhost:9093) which it was generated after inserting row in database


