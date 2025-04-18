#!/bin/bash
curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "yellow"
curl -X POST http://localhost:17000 -d "figure 0.3 0.3"
curl -X POST http://localhost:17000 -d "update"

for i in {1..20}
do
  curl -X POST http://localhost:17000 -d "move 0.01 0.01"
  curl -X POST http://localhost:17000 -d "update"
  sleep 1
done
