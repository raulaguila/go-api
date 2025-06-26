# Go API Template

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

### Summary

1. [Description](#description-)
2. [Prerequisites](#prerequisites-)
3. [Makefile](#makefile-)
4. [Starting the Project](#starting-the-project-)
5. [Stopping the Project](#stopping-the-project-)
6. [Features](#features-)
    1. [Profile Module](#profile-module-http)
    2. [User Module](#user-module-http)
    3. [Authentication Module](#authentication-module-http)
    4. [Product Module](#product-module-http)
7. [Code Status](#code-status-)
8. [Contributors](#contributors-)
9. [License](#license-)

<h1></h1>

1. #### Description [&uarr;](#summary)

   User-friendly API template solution designed as a foundation for more complex APIs.

2. #### Prerequisites [&uarr;](#summary)
    * Docker
    * Docker Compose
    * Golang 1.24+ (Optional)

3. #### Makefile [&uarr;](#summary)
   <details>
   <summary>Makefile commands:</summary>

    ```shell
    Usage:
        make [COMMAND]
        make help

    Commands:
    
    init                           Create environment file
    build                          Build the application from source code
    run                            Run application from source code
    compose-build-services         Create and start services containers
    compose-build-built            Create and start containers from built
    compose-build-source           Create and start containers from source code
    compose-down                   Stop and remove containers and networks
    compose-remove                 Stop and remove containers, networks and volumes
    compose-exec                   Access container bash
    compose-log                    Show container logger
    compose-top                    Display containers processes
    compose-stats                  Display containers stats
    go-test                        Run test and generate coverage report
    go-benchmark                   Benchmark code performance
    go-lint                        Run lint checks
    go-audit                       Conduct quality checks
    go-swag                        Update swagger files
    go-format                      Fix code format issues
    go-tidy                        Clean and tidy dependencies
    ```
   </details>

4. #### Starting the Project [&uarr;](#summary)
    * Download and extract the latest build [release](https://github.com/raulaguila/go-api/releases)
    * Open the terminal in the release folder
    * Run `make compose-build`
    * Or to build from source code: `make compose-build COMPOSE=build/source.compose.yml`

5. #### Stopping the Project [&uarr;](#summary)
    * Open the terminal in the release folder
    * Run `make compose-remove`

6. #### Features [&uarr;](#summary)

    * Get default user email and password on environment file `configs/.env`
    * Test API endpoints accessing [swagger page](http://127.0.0.1:9000/swagger)

    1. ###### Profile Module ([HTTP](../api/profile.http))

       | Endpoint        | HTTP Method |       Description        |
       |:----------------|:-----------:|:------------------------:|
       | `/profile`      |    `GET`    |    `Get all profiles`    |
       | `/profile`      |   `POST`    |   `Insert new profile`   |
       | `/profile`      |  `DELETE`   | `Delete profiles by IDs` |
       | `/profile/{id}` |    `GET`    |   `Get profile by ID`    |
       | `/profile/{id}` |    `PUT`    |  `Update profile by ID`  |

    2. ###### User Module ([HTTP](../api/user.http))

       | Endpoint     | HTTP Method |       Description       |
       |:-------------|:-----------:|:-----------------------:|
       | `/user`      |    `GET`    |     `Get all users`     |
       | `/user`      |   `POST`    |      `Insert user`      |
       | `/user`      |  `DELETE`   |      `Delete user`      |
       | `/user/{id}` |    `GET`    |    `Get user by ID`     |
       | `/user/{id}` |    `PUT`    |   `Update user by ID`   |
       | `/user/pass` |    `PUT`    |  `Set user's password`  |
       | `/user/pass` |  `DELETE`   | `Reset user's password` |

    3. ###### Authentication Module ([HTTP](../api/auth.http)):

       | Endpoint | HTTP Method |               Description               |
       |:---------|:-----------:|:---------------------------------------:|
       | `/auth`  |   `POST`    |          `User authentication`          |
       | `/auth`  |    `GET`    |  `User authenticated via access token`  |
       | `/auth`  |    `PUT`    | `User refresh tokens via refresh token` |

        * Pass token using prefix _**Bearer**_ in Authorization request header:

       ```bash
       Authorization: Bearer <token>
       ```

    4. ###### Product Module ([HTTP](../api/product.http))

       | Endpoint        | HTTP Method |      Description       |
       |:----------------|:-----------:|:----------------------:|
       | `/product`      |    `GET`    |   `Get all products`   |
       | `/product`      |   `POST`    | `Insert a new product` |
       | `/product`      |  `DELETE`   |   `Delete products`    |
       | `/product/{id}` |    `GET`    |  `Get product by ID`   |
       | `/product/{id}` |    `PUT`    | `Update product by ID` |

7. #### Code Status [&uarr;](#summary)
    * Development

8. #### Contributors [&uarr;](#summary)

   <a href="https://github.com/raulaguila" target="_blank">
     <img src="https://contrib.rocks/image?repo=raulaguila/go-api" alt="raulaguila">
   </a>

9. #### License [&uarr;](#summary)

   Copyright © 2023 [raulaguila](https://github.com/raulaguila). This project is [MIT](../LICENSE) licensed.
