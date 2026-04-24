#!/usr/bin/env bash 

sqlite3 ~/.config/taskmanager/tasks.db ".backup db/backups/taskmanager-$(date +%Y-%m-%d-%H:%M:%S).db"
