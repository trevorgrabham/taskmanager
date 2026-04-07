#!/usr/bin/env bash 

PID_NAME='/tmp/taskmanager-webserver.pid'
parent_process="$(cat $PID_NAME 2>/dev/null)"

if ps --pid "$(cat /tmp/taskmanager-webserver.pid 2>/dev/null)" > /dev/null 2>&1; then 
  echo "Server Running..."
  echo "Stopping Server..."
  children_processes="$(pgrep -P $parent_process)"
  if [ -n "$children_processes" ]; then
    kill $children_processes $parent_process
  else
    kill $parent_process
  fi
  echo "Server Stopped"
  echo "Generating Templates..."
  templ generate
  echo "Templates Generated"
  echo "Compiling CSS..."
  cat ./static/*/index.css | csso > ./static/index.min.css
  echo "CSS Compiled"
  echo "Restarting Server..."
  go run ./cmd/web & 
  echo "$!" > "$PID_NAME"
  echo "Server Running..."
else 
  echo "Generating Templates..."
  templ generate
  echo "Templates Generated"
  echo "Compiling CSS..."
  cat ./static/*/index.css | csso > ./static/index.min.css
  echo "CSS Compiled"
  echo "Starting Server..."
  go run ./cmd/web & 
  echo "$!" > "$PID_NAME"
  echo "Server Running..."
fi 
