#!/usr/bin/env bash

PID_NAME='/tmp/taskmanager-webserver.pid'
parent_process="$(cat $PID_NAME)"

if ps --pid "$(cat /tmp/taskmanager-webserver.pid)" > /dev/null 2>&1; then 
  echo "Stopping Server..."
  children_processes="$(pgrep -P $parent_process)"
  if [ -n "$children_processes" ]; then
    kill $children_processes $parent_process
  else
    kill $parent_process
  fi
  echo "Server Stopped"
fi
