#!/bin/bash
#start-ganache-and-save-keyssh

# Debug output
set -x

# Start Ganache and capture output
ganache --mnemonic "PVote" > ganache_output.txt &

# Wait for Ganache
sleep 5

# List Node.js processes
echo "Before killing process:"
pgrep -fl node

# Find the process on port 8545
GANACHE_PID=$(lsof -t -i:8545)
if [ ! -z "$GANACHE_PID" ]; then
  echo "Killing process $GANACHE_PID occupying port 8545"
  sudo kill $GANACHE_PID
  sleep 2
else
  echo "No process found occupying port 8545"
fi

# Check process termination
if sudo kill -0 $GANACHE_PID 2>/dev/null; then
  echo "Failed to kill process $GANACHE_PID"
else
  echo "Process $GANACHE_PID killed successfully"
fi

# Initialize .env
echo "" > .env

# Write keys without the 0x prefix to .env
cat ganache_output.txt | grep 'Private Keys' -A 12 | grep -o '0x.*' | while read -r line; do
  echo "PRIVATE_KEY_$((++i))=${line:2}" >> .env
done

# Remove temporary files
rm ganache_output.txt

echo "Private keys saved to .env file."
