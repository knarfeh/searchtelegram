#!/usr/bin/env sh

sync
chmod +x /*.sh
sync
supervisord --nodaemon --configuration /etc/supervisord.conf
