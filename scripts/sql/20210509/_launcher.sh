#!/usr/bin/sh
cat 000_alter_deleted_cols.sql | docker exec -i diosteama_db_1 psql -U diosteama