#!/bin/bash
set -e
DB_NAME=${1:-bank}
DB_USER=${2:-bank_app_user}
DB_USER_PASS=${3:-Pwd2022}
sudo su postgres <<EOF
createdb $DB_NAME;
createdb test_$DB_NAME;
psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_USER_PASS';"
psql -c "grant all privileges on database $DB_NAME to $DB_USER;"
psql -c "grant all privileges on database test_$DB_NAME to $DB_USER;"
echo "Postgres User '$DB_USER' and databases '$DB_NAME', 'test_$DB_NAME' created."
EOF