#!/bin/bash

declare -A tokens
for ((i = 0; i < 2; i++)); do
    openssl genrsa -out pvt.pem 3072
    openssl rsa -in pvt.pem -pubout -out pub.pem >/dev/null 2>&1

    # shellcheck disable=SC2002
    tokens[$i, 0]=$(cat pvt.pem | base64 | tr -d \\n)
    # shellcheck disable=SC2002
    tokens[$i, 1]=$(cat pub.pem | base64 | tr -d \\n)

    rm pvt.pem pub.pem
done

echo "TZ='America/Manaus'                             # Set system time zone
SYS_PREFORK='1'                                 # Enable Fiber Prefork

API_PORT='9000'                                 # API Container PORT
API_LOGGER='1'                                  # API Logger enable
API_SWAGGO='1'                                  # API Swagger enable
API_DEFAULT_SORT='updated_at'                   # API default column sort
API_DEFAULT_ORDER='desc'                        # API default order
API_ACCEPT_SKIP_AUTH='1'                        # API accept skip auth header

ACCESS_TOKEN_EXPIRE='15'                        # [MINUTES] Access token expiration time
ACCESS_TOKEN_PRIVAT='${tokens[0, 0]}'           # Token to encode access token - PRIVATE TOKEN
ACCESS_TOKEN_PUBLIC='${tokens[0, 1]}'           # Token to decode access token - PUBLIC TOKEN

RFRESH_TOKEN_EXPIRE='60'                        # [MINUTES] Refresh token expiration time
RFRESH_TOKEN_PRIVAT='${tokens[1, 0]}'           # Token to encode refresh token - PRIVATE TOKEN
RFRESH_TOKEN_PUBLIC='${tokens[1, 1]}'           # Token to decode refresh token - PUBLIC TOKEN

POSTGRES_HOST='postgres'                        # Postgres Container HOST
POSTGRES_PORT='5432'                            # Postgres Container PORT
POSTGRES_USER='root'                            # Postgres USER
POSTGRES_PASS='root'                            # Postgres PASS
POSTGRES_BASE='go_api'                          # Postgres BASE
" >.env
