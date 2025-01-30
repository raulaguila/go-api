#!/bin/bash

access_token=$(openssl genrsa 2048 | base64 | tr -d \\n)
refresh_token=$(openssl genrsa 2048 | base64 | tr -d \\n)

echo "TZ='America/Manaus'                             # Set system time zone

API_PORT='9000'                                 # API Container PORT
API_LOGGER='1'                                  # API Logger enable
API_SWAGGO='1'                                  # API Swagger enable
API_ENABLE_PREFORK='1'                          # API enable fiber prefork
API_DEFAULT_SORT='updated_at'                   # API default column sort
API_DEFAULT_ORDER='desc'                        # API default order
API_ACCEPT_SKIP_AUTH='1'                        # API accept skip auth header

ACCESS_TOKEN_EXPIRE='15'                        # Access token expiration time in minutes
RFRESH_TOKEN_EXPIRE='60'                        # Refresh token expiration time in minutes

ACCESS_TOKEN_PRIVAT='${access_token}'           # Token to encode access token - PRIVATE TOKEN
RFRESH_TOKEN_PRIVAT='${refresh_token}'          # Token to encode refresh token - PRIVATE TOKEN

POSTGRES_HOST='postgres'                        # Postgres Container HOST
POSTGRES_PORT='5432'                            # Postgres Container PORT
POSTGRES_USER='root'                            # Postgres USER
POSTGRES_PASS='root'                            # Postgres PASS
POSTGRES_BASE='go_api'                          # Postgres BASE
" >.env
