#!/bin/sh
 
DATABASE=todos.db
 
cat "db_init.sql" | sqlite3 "${DATABASE}"
