#!/usr/bin/env sh

sync
chmod +x /*.sh /bin/searchtelegram
/bin/searchtelegram download_cert

mv searchtelegramdotcom.key searchtelegramdotcom_bundle.crt /etc/nginx/ssl/

sync
supervisord --nodaemon --configuration /etc/supervisord.conf
