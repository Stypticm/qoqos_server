#!/bin/bash
echo "--- DB ENCODING ---"
psql -U qoqos -d postgres -c "SHOW server_encoding; SHOW client_encoding;"
echo "--- SAMPLE DATA (HEX) ---"
psql -U qoqos -d qoqos_db -t -c "SELECT title, encode(title::bytea, 'hex') FROM \"BlogPost\" LIMIT 1;"
echo "--- TABLE LIST ---"
psql -U qoqos -d qoqos_db -c "\dt" | grep -E "BlogPost|MarketplaceLot"
