#!/usr/bin/env sh

sync
chmod +x /*.sh /bin/searchtelegram

sync
supervisord --nodaemon --configuration /etc/supervisord.conf
