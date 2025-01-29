#!/bin/bash

DB_FILE="db/wiki.db"  # Path to your SQLite database file

if [ ! -f "$DB_FILE" ]; then
    echo "Error: Database file '$DB_FILE' not found!"
    exit 1
fi

echo "Database Schema for $DB_FILE:"
sqlite3 "$DB_FILE" ".schema"

