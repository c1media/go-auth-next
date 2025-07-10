#!/bin/bash

# Function to load env file
load_env() {
  if [ -f .env.local ]; then
    export $(cat .env.local | sed 's/#.*//g' | xargs)
  elif [ -f .env ]; then
    export $(cat .env | sed 's/#.*//g' | xargs)
  fi
}

# Load environment variables
load_env

# Start the back-end server with env loaded
cd authserver
export $(cat ../.env.local | sed 's/#.*//g' | xargs) && go run ./cmd/server &
BACKEND_PID=$!
cd ..

# Start the front-end server
cd front-end
PORT=3000 npm run dev &
FRONTEND_PID=$!

# Wait for both processes
wait $BACKEND_PID $FRONTEND_PID
