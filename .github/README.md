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

#### Summary

1. [Description](#description-)
2. [Dependencies](#dependencies-)
3. [Make Command](#make-file-)
4. [Running the project](#running-the-project-)
   1. [From repository clone](#from-repository-clone)
   2. [From built version](#from-built-version)
5. [Stopping the project](#stopping-the-project-)

<h1></h1>

1. #### Description [&uarr;](#summary)

   User-friendly API template solution designed as a foundation for more complex APIs.

2. #### Dependencies [&uarr;](#summary)

   - Make
   - Docker
   - Docker Compose
   - Go 1.24+ (Optional)

3. #### Make File [&uarr;](#summary)

      <details>
      <summary>Commands:</summary>

   ```sh
   Usage:
   make [COMMAND]

   Example:
   make build

   Commands:

   help                           Display available commands and their descriptions
   init                           Create environment file
   test                           Run tests and generate coverage report
   run                            Run application from source code
   build                          Build the all applications from source code
   swag                           Update swagger files
   format                         Fix code format issues
   tidy                           Clean and tidy dependencies
   lint                           Run lint checks
   audit                          Conduct quality checks
   benchmark                      Benchmark code performance
   compose-up                     Create and start containers
   compose-build                  Build, create and start containers
   compose-down                   Stop and remove containers and networks
   compose-clean                  Clear dangling Docker images
   compose-remove                 Stop and remove containers, networks and volumes
   compose-exec                   Access container bash
   compose-log                    Show container logger
   compose-top                    Display containers processes
   compose-stats                  Display containers stats
   ```

      </details>

4. #### Running the project [&uarr;](#summary)

   If you have the cloned repository, follow the instructions in the item "[From repository clone](#from-repository-clone)"

   If you downloaded the released version, follow the instructions in the item "[From released version](#from-released-version)"

   1. ##### From repository clone

      1. Open the terminal in the cloned repository folder.
         - If there is a golang installation on the system:
           - Run `make init build compose-build` to create environment file, build the application and create and start containers.
         - If there is no golang installation on the system:
           - Run `make init compose-build base=source` to create environment file and create and start production containers.

   2. ##### From released version

      1. Open terminal in built version folder.
      2. Run `make compose-build` to create and start production containers.

5. #### Stopping the project [&uarr;](#summary)

   1. Open the terminal in the project folder.
   2. Run `make compose-down` to stop and remove containers and networks or run `make compose-remove` to stop and
      remove containers, networks and volumes.

<p style="text-align:right">&#40;<a href="#title">back to top</a>&#41;</p>


6. #### Features [&uarr;](#summary)

    * Get default user email and password on environment file `configs/.env`
    * Test API endpoints using [http files](../api) or accessing [swagger page](http://127.0.0.1:9000/swagger)

    1. ###### Profile Module

       | Endpoint        | HTTP Method |       Description        |
       |:----------------|:-----------:|:------------------------:|
       | `/profile`      |    `GET`    |    `Get all profiles`    |
       | `/profile`      |   `POST`    |   `Insert new profile`   |
       | `/profile`      |  `DELETE`   | `Delete profiles by IDs` |
       | `/profile/{id}` |    `GET`    |   `Get profile by ID`    |
       | `/profile/{id}` |    `PUT`    |  `Update profile by ID`  |

    2. ###### User Module

       | Endpoint     | HTTP Method |       Description       |
       |:-------------|:-----------:|:-----------------------:|
       | `/user`      |    `GET`    |     `Get all users`     |
       | `/user`      |   `POST`    |      `Insert user`      |
       | `/user`      |  `DELETE`   |      `Delete user`      |
       | `/user/{id}` |    `GET`    |    `Get user by ID`     |
       | `/user/{id}` |    `PUT`    |   `Update user by ID`   |
       | `/user/pass` |    `PUT`    |  `Set user's password`  |
       | `/user/pass` |  `DELETE`   | `Reset user's password` |

    3. ###### Authentication Module

       | Endpoint | HTTP Method |               Description               |
       |:---------|:-----------:|:---------------------------------------:|
       | `/auth`  |   `POST`    |          `User authentication`          |
       | `/auth`  |    `GET`    |  `User authenticated via access token`  |
       | `/auth`  |    `PUT`    | `User refresh tokens via refresh token` |

        * Pass token using prefix _**Bearer**_ in Authorization request header:

       ```bash
       Authorization: Bearer <token>
       ```

7. #### Code Status [&uarr;](#summary)
    * Development

8. #### Contributors [&uarr;](#summary)

   <a href="https://github.com/raulaguila" target="_blank">
     <img src="https://contrib.rocks/image?repo=raulaguila/go-api" alt="raulaguila">
   </a>

9. #### License [&uarr;](#summary)

   Copyright © 2023 [raulaguila](https://github.com/raulaguila). This project is [MIT](../LICENSE) licensed.