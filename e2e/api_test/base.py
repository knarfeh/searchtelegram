#!/usr/bin/env python
# -*- coding: utf-8 -*-

import logging

import redis
from maic.testrunner import TestRunner
from maic.assertion import Assertion
from . import configs

LOGGER = logging.getLogger()


class TestBase(Assertion):
    """
    """
    def __init__(self):
        r = redis.StrictRedis(host=configs.REDIS_HOST, port=configs.REDIS_PORT, db=0)
        self.r = r
    def assert_successful(self, response):
        msg = response.get('message', response.get('total', response))
        LOGGER.info("Assert successful, msg: %s", msg)
        self.assert_true(response['success'], msg=msg)
        return response
