#!/bin/bash
psql -U qoqos -d qoqos_db -c "SELECT title, encode(title::bytea, 'hex') FROM \"BlogPost\" LIMIT 3;"
