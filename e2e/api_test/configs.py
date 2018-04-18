#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import logging
import json

__author__ = "knarfeh@outlook.com"


def str2bool(v):
    return v and v.lower() in ("yes", "true", "t", "1")


def get_list_from_str(string, separator=','):
    if string is not None and string != '':
        return string.split(separator)


def get_float_from_str(string, default=0):
    try:
        return float(string)
    except:
        return default


def get_json_from_str(string, default={"testint": "1.1.1.1"}):
    try:
        return json.loads(string)
    except:
        return default

def get_header():
    headers = {
        "content-type": "application/json"
    }
    return headers

RECIPIENTS = get_list_from_str(os.getenv('RECIPIENTS', 'knarfeh@outlook.com'))
EMAIL = {
    'from': 'knarfeh@outlook.com',
    'recipients': RECIPIENTS
}

RETRY_TIMES = 1
ENV = os.environ.get('ENV') or "STAGING"

API_URL = os.getenv("STAPIURL", "http://192.168.199.121:5000")
REDIS_HOST = os.getenv("REDISHOST", "localhost")
REDIS_PORT = os.getenv("REDISPORT", 16379)

# test cases list
TESTCASES = get_list_from_str(os.getenv("TESTCASES"))


EMAIL_CONFIG = {
    'template_path': 'email',
    'sender': os.getenv('EMAIL_FROM', '2559775198@qq.com'),
    'recipient': 'hejun1874@gmail.com'
}
MAILGUN_CONFIG = {
    'MAIL_DOMAIN_NAME': os.getenv('MAIL_DOMAIN_NAME', None),
    'MAIL_API_KEY': os.getenv('MAIL_API_KEY', None)
}


logging.basicConfig(
    level=logging.DEBUG,
    format=('%(asctime)s [%(process)d] %(levelname)s %(pathname)s' +
            ' %(funcName)s Line:%(lineno)d %(message)s'),
)

