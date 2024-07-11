<h1 style="text-align:center">Go API Template</h1>

<p style="text-align:center">
  <a href="https://github.com/raulaguila/go-api/releases" target="_blank" style="text-decoration: none;">
    <img src="https://img.shields.io/github/v/release/raulaguila/go-api.svg?style=flat&labelColor=0D1117" alt="release">
  </a>
  <img src="https://img.shields.io/github/repo-size/raulaguila/go-api?style=flat&labelColor=0D1117" alt="size">
  <img src="https://img.shields.io/github/stars/raulaguila/go-api?style=flat&labelColor=0D1117" alt="stars">
  <a href="../LICENSE" target="_blank" style="text-decoration: none;">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat&labelColor=0D1117" alt="license">
  </a>
  <a href="https://goreportcard.com/report/github.com/raulaguila/go-api" target="_blank" style="text-decoration: none;">
    <img src="https://goreportcard.com/badge/github.com/raulaguila/go-api?style=flat&labelColor=0D1117" alt="report">
  </a>
  <a href="https://github.com/raulaguila/go-api/actions?query=workflow%3Ago-test" target="_blank" style="text-decoration: none;">
    <img src="https://github.com/raulaguila/go-api/actions/workflows/go_test.yml/badge.svg" alt="test">
  </a>
  <a href="https://github.com/raulaguila/go-api/actions?query=workflow%3Ago-build" target="_blank" style="text-decoration: none;">
    <img src="https://github.com/raulaguila/go-api/actions/workflows/go_build.yml/badge.svg" alt="build">
  </a>
</p>

## Prerequisites

- Docker

## Getting Started

- Help with make command

```sh
Usage:
        make [COMMAND]
        make help

Commands: 

init                           Create environment file
compose-up                     Run 'docker compose up -d' to create and start containers
compose-build-dev              Run 'docker compose --profile dev up -d --build' to create and start containers
compose-build-prod             Run 'docker compose --profile prod up -d --build' to create and start containers
compose-down                   Run 'docker compose --profile all down' to stop and remove containers and networks
compose-remove                 Run 'docker compose --profile all down -v --remove-orphans' to stop and remove containers, networks and volumes
compose-exec                   Run 'docker compose exec -it backend bash' to access container bash
compose-log                    Run 'docker compose logs -f backend' to show container logger
compose-top                    Run 'docker compose top' to display containers processes
compose-stats                  Run 'docker compose stats' to display containers stats
go-run                         Run application from source code
go-test                        Run tests and generate coverage report
go-build                       Build the application from source code
go-benchmark                   Benchmark code performance
go-lint                        Run lint checks
go-audit                       Conduct quality checks
go-swag                        Update swagger files
go-format                      Fix code format issues
go-tidy                        Clean and tidy dependencies
```

- Run project

1. Download and extract the latest build [release](https://github.com/raulaguila/go-api/releases)
2. Open the terminal in the release folder
3. Run:

```sh
make compose-build
```

- Remove project

```sh
make compose-remove
```

## Features

- Get default user email and password on environment file `configs/.env`
- Test API endpoints using <a href="../api" target="_blank">http files</a> or
  accessing <a href="http://127.0.0.1:9000/swagger/index.html" target="_blank">swagger page</a>

[Profile module](../api/profile.http):

| Endpoint        | HTTP Method |       Description        |
|:----------------|:-----------:|:------------------------:|
| `/profile`      |    `GET`    |    `Get all profiles`    |
| `/profile`      |   `POST`    |   `Insert new profile`   |
| `/profile`      |  `DELETE`   | `Delete profiles by IDs` |
| `/profile/{id}` |    `GET`    |   `Get profile by ID`    |
| `/profile/{id}` |    `PUT`    |  `Update profile by ID`  |

[User module](../api/user.http):

| Endpoint           | HTTP Method |       Description       |
|:-------------------|:-----------:|:-----------------------:|
| `/user`            |    `GET`    |     `Get all users`     |
| `/user`            |   `POST`    |      `Insert user`      |
| `/user`            |  `DELETE`   |      `Delete user`      |
| `/user/{id}`       |    `GET`    |    `Get user by ID`     |
| `/user/{id}`       |    `PUT`    |   `Update user by ID`   |
| `/user/{id}/photo` |    `GET`    |   `Get user's photo`    |
| `/user/{id}/photo` |    `PUT`    |   `Set user's photo`    |
| `/user/pass`       |    `PUT`    |  `Set user's password`  |
| `/user/pass`       |  `DELETE`   | `Reset user's password` |

[Authentication module](../api/auth.http):

| Endpoint | HTTP Method |               Description               |
|:---------|:-----------:|:---------------------------------------:|
| `/auth`  |   `POST`    |          `User authentication`          |
| `/auth`  |    `GET`    |  `User authenticated via access token`  |
| `/auth`  |    `PUT`    | `User refresh tokens via refresh token` |

- Pass token using prefix _**Bearer**_ in Authorization request header:

```bash
Authorization: Bearer <token>
```

[Department module](../api/department.http):

| Endpoint           | HTTP Method |        Description        |
|:-------------------|:-----------:|:-------------------------:|
| `/department`      |    `GET`    |   `Get all departments`   |
| `/department`      |   `POST`    |  `Insert new department`  |
| `/department`      |  `DELETE`   |   `Delete departments`    |
| `/department/{id}` |    `GET`    |  `Get department by ID`   |
| `/department/{id}` |    `PUT`    | `Update department by ID` |

## Code status

- Development

## Contributors

<a href="https://github.com/raulaguila" target="_blank">
  <img src="https://contrib.rocks/image?repo=raulaguila/go-api" alt="raulaguila">
</a>

## License

Copyright © 2023 [raulaguila](https://github.com/raulaguila).
This project is [MIT](../LICENSE) licensed.
