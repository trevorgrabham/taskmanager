#!/usr/bin/env bash 

sqlite3 ~/.config/taskmanager/tasks.db ".backup internal/db/backups/taskmanager_$(date +%Y%m%d_%H%M%S).db"
