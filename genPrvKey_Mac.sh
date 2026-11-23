#!/bin/bash
#start-ganache-and-save-keyssh

# Debug output
set -x

# Start Ganache and capture output
OS=$(uname -s)
 
case "$OS" in
  Linux*)
    echo "Linux"
    ganache --mnemonic "PVote"  > ganache_output.txt &
    ;;
  Darwin*)
    echo "macOS"
    ganache-cli --mnemonic "PVote" > ganache_output.txt &
    ;;
  CYGWIN*|MINGW32*|MSYS*|MINGW*)
    echo "Windows"
    ;;
  *)
    echo "Unknown OS"
    ;;
esac

# Wait for Ganache
sleep 5
rm .env

# Write accounts to .env
i=1
cat ganache_output.txt | grep -A 12 'Available Accounts' | grep '0x' | while read -r line; do
  address=$(echo $line | awk '{print $2}')
  echo "ACCOUNT_$i=$address" >> .env
  ((i++))
done
a=0
# Write keys without the 0x prefix to .env
cat ganache_output.txt | grep 'Private Keys' -A 12 | grep -o '0x.*' | while read -r line; do
  echo "PRIVATE_KEY_$((++a))=${line:2}" >> .env
done
#On Ubuntu, Ganache may require manual cleanup
rm ganache_output.txt
ps -ef|grep 'ganache'|xargs kill -9
