#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = "knarfeh@outlook.com"

import logging

from maic.testrunner import TestRunner

from api_test.api_tgbot import TgBotTest
from api_test import config


class ApiTest(TestRunner, TgBotTest):
    data = {}

    def set_up(self):
        print("WTF")
        self.data = {}
        self.tgbot_set_up()

    def tear_down(self):
        print("WTF???")
        self.data = {}
        self.tgbot_tear_down()
