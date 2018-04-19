#!/usr/bin/env sh

sync
chmod +x /src/*.sh
sync
python main.py
supervisord --nodaemon --configuration /etc/supervisord.conf
