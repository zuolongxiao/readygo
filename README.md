# ReadyGo

> A Golang framework for quick API development

## Features

- Based on Gin + GORM
- Follow RESTful API conventions
- JWT authentication
- RBAC permission control
- High performance

## Quick start

> Before started, you need to have MySQL and Redis server installed if you are intended to use MySQL or Redis, MySQL and Redis is optional.

1. Pull code

    ```bash
    git clone https://github.com/zuolongxiao/readygo.git
    ```

2. Download and install dependencies

    ```bash
    # For users in China, optional
    go env -w GOPROXY=https://goproxy.cn,direct

    cd readygo
    go mod tidy

    # Build, will generate an executable file named `readygo`
    go build
    ```

3. Configuration

    > You can use config file or ENV variables as configuration options, for more details, please see `config.sample.yaml`.
    > `readygo` will read `config.yaml` in the same directory if present.

    ```bash
    # using config file
    cp config.sample.yaml config.yaml
    vi config.yaml
    # change the following options:
    # JWT.Secret: <your secret>
    # Database.Type: MySQL # MySQL or SQLite

    # if Database.Type is MySQL
    # Database.Host: 127.0.0.1
    # Database.Port: 3306
    # Database.User: root
    # Database.Password: <your MySQL server password>

    # using ENV variables
    export JWT_SECRET=<your secret>
    export  DATABASE_TYPE=MySQL # MySQL or SQLite

    # if DATABASE_TYPE is MySQL
    export DATABASE_HOST=127.0.0.1
    export DATABASE_PORT=3306
    export DATABASE_USER=root
    export DATABASE_PASSWORD=<your MySQL server password>
    ```

4. Create database (if you are using MySQL)

    ```bash
    mysql -h127.0.0.1 -P3306 -uroot -p -e "create database readygo"
    ```

5. Initialization

    ```bash
    # Migrate tables
    ./readygo admin migrate

    # Load permissions into database, optional
    ./readygo admin permission

    # Create first administrator, as default super administrator
    ./readygo admin create -u admin -p <your password>
    ```

6. Starting HTTP service

    ```bash
    ./readygo serve
    ```

7. Testing

    ```bash
    # Obtain JWT token, this token need to be carried in the future request
    curl --request POST 'http://127.0.0.1:9331/api/auth' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "username": "admin",
        "password": "<your password>"
    }'

    # Fetch user list
    curl --request GET 'http://127.0.0.1:9331/api/v1/admins' \
    --header 'Authorization: Bearer <JWT TOKEN>'
    ```
