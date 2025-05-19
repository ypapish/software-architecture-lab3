#!/bin/bash

# Скидання стану
curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "figure 0.3 0.3"
curl -X POST http://localhost:17000 -d "update"

for i in {1..7}
do
  curl -X POST http://localhost:17000 -d "move 0.07 0.07"
  curl -X POST http://localhost:17000 -d "update"
  sleep 0.01
done
