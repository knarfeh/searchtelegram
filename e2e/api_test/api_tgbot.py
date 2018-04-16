#!/usr/bin/env python
# -*- coding: utf-8 -*-


import logging
import json
from time import time

import requests

from maic.decorators import case
from api_test import config
from api_test.base import TestBase

__author__ = "knarfeh@outlook.com"
LOGGER = logging.getLogger()


def get_last_message():
    pass


class TgBotTest(TestBase):
    """
    """

    def tgbot_set_up(self):
        LOGGER.debug("TODO: tgbot set up")

    def tgbot_tear_down(self):
        LOGGER.debug("TODO: tgbot tear down")

    @case
    def tgbot_test__ping(self):
        LOGGER.debug("test ping")
        print("test ping")
        result = {
            "success": False,
            "message": "message",
        }
        self.assert_successful(result)
        return
