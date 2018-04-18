#!/usr/bin/env python
# -*- coding: utf-8 -*-


import logging
import json
from time import time, sleep

import requests

from maic.decorators import case
from api_test import configs
from api_test.base import TestBase

__author__ = "knarfeh@outlook.com"
LOGGER = logging.getLogger()


def get_tgbot_response(payload, header):
    LOGGER.info("tgbot ping")
    start_time = time()
    url = "{}/api/tgbot".format(configs.API_URL)
    LOGGER.debug("URL: %s", url)
    r = requests.post(url, headers=header, data=json.dumps(payload))
    end_time = time()
    response = {
        "status": r.status_code,
        "text": r.text,
        "total": end_time - start_time
    }
    return response


class TgBotTest(TestBase):
    """
    """

    def tgbot_set_up(self):
        LOGGER.debug("Set up tgbot test ...")
        payload = self._get_tgbot_payload("/echo Start testing tgbot ...")
        get_tgbot_response(payload, configs.get_header())
        sleep(5)

    def tgbot_tear_down(self):
        LOGGER.debug("tgbot tear down")
        payload = self._get_tgbot_payload("/echo End testing tgbot")
        get_tgbot_response(payload, configs.get_header())

    def _get_last_message(self):
        pass

    def _get_tgbot_payload(self, text):
        payload = {
            "message": {
                "message_id": 204,
                "from": {
                "id": 312172714,
                "first_name": "knarfeh",
                "last_name": "",
                "username": "knarfeh"
                },
                "date": 1522474188,
                "chat": {
                "id": 312172714,
                "type": "private",
                "title": "",
                "first_name": "knarfeh",
                "last_name": "",
                "username": "knarfeh"
                },
                "text": text
            }
        }
        return payload

    @case
    def tgbot_test__othercommand(self):
        payload = self._get_tgbot_payload("/delete knarfeh")
        get_tgbot_response(payload, configs.get_header())

        print("Sleep 5 seconds \n")
        sleep(5)
        print("Tesing submit \n")
        self.test_submit(exist=False)

        print("Sleep 5 seconds \n")
        sleep(5)
        print("Tesing ping \n")
        self.test_ping()

        print("Sleep 5 seconds \n")
        sleep(5)
        print("Tesing echo \n")
        self.test_echo()

        print("Sleep 5 seconds \n")
        sleep(5)
        print("Tesing start \n")
        self.test_start()

        print("Sleep 5 seconds \n")
        sleep(5)
        print("Tesing get \n")
        self.test_get()

        print("Sleep 5 seconds \n")
        sleep(5)
        print("Tesing submit exist \n")
        self.test_submit(exist=True)

        print("Sleep 5 seconds \n")
        sleep(5)
        print("Tesing stats \n")
        self.test_stats()

        print("Sleep 10 seconds \n")
        sleep(10)
        print("Tesing delete \n")
        self.test_delete()

        print("Sleep 3 seconds")
        sleep(3)
        print("Tesing search \n")
        self.test_search()

    def test_ping(self):
        LOGGER.debug("test ping command")
        payload = self._get_tgbot_payload("/ping")
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot ping failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            if last_message == "pong ":
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "ping command, got last message: {}, should be: {}".format(last_message, "pong")
                }
        self.assert_successful(result)
        return

    def test_echo(self):
        LOGGER.debug("Test echo command")
        echo_message = "Use echo command to test echo"
        payload = self._get_tgbot_payload("/echo {}".format(echo_message))
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot echo command failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            if last_message == echo_message:
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "echo command, got last message: {}, should be: {}".format(last_message, echo_message)
                }
        self.assert_successful(result)
        return

    def test_start(self):
        LOGGER.debug("Test start command")
        payload = self._get_tgbot_payload("/start")
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot start command failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            LOGGER.debug("Start command, got message: %s", last_message)
            if "I will help you search" in last_message:
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "echo command, got last message: {}".format(last_message)
                }
        self.assert_successful(result)
        return

    def test_get(self):
        LOGGER.debug("Test get command")
        payload = self._get_tgbot_payload("/get knarfeh")
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot get command, failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            sleep(3)
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            LOGGER.debug("Start command, got message: %s", last_message)
            if "@knarfeh" in last_message:
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "get command, got last message: {}".format(last_message)
                }
        self.assert_successful(result)
        return

    def test_stats(self):
        LOGGER.debug("Test stats command")
        payload = self._get_tgbot_payload("/stats")
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot get command, failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            sleep(5)
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            LOGGER.debug("stats command, got message: %s", last_message)
            if "Unique user" in last_message:
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "stats command, got last message: {}".format(last_message)
                }
        self.assert_successful(result)
        return

    def test_delete(self):
        LOGGER.debug("Test delete command")
        payload = self._get_tgbot_payload("/delete knarfeh")
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot delete command, failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            result = {
                "success": True,
                "total": tgbot_result["total"]
            }
        self.assert_successful(result)

        print("Wait for redis delete pipeline, sleep 5 seconds")
        sleep(5)
        payload = self._get_tgbot_payload("/get knarfeh")
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot get command, failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            sleep(10)
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            LOGGER.debug("After delete command, try test command, get last_message: %s", last_message)
            if "this id does not exist" in last_message:
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "After delete, try get command, got last message: {}".format(last_message)
                }
        self.assert_successful(result)
        return

    def test_submit(self, exist=True):
        LOGGER.debug("Test submit")
        payload = self._get_tgbot_payload("/submit knarfeh")
        tgbot_result = get_tgbot_response(payload, configs.get_header())
        if exist is True:
            should_message = "this id already exist"
        else:
            should_message = "Successfully submitted"
        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot submit command, failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            sleep(10)
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            LOGGER.debug("submit command, got message: %s", last_message)
            if should_message in last_message:
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "submit command, got last message: {}".format(last_message)
                }
        self.assert_successful(result)
        return

    def test_search(self):
        LOGGER.debug("Test search")
        payload = self._get_tgbot_payload("/search telegram")
        tgbot_result = get_tgbot_response(payload, configs.get_header())

        if tgbot_result["status"] != 200:
            result = {
                "success": False,
                "message": "tgbot submit command, failed, error code: {}, error message: {}".format(
                    tgbot_result["status"],
                    tgbot_result["text"]
                )
            }
        else:
            sleep(5)
            last_message = self.r.get("e2e:last-message").decode("utf-8")
            LOGGER.debug("submit command, got message: %s", last_message)
            if "@telegram" in last_message:
                result = {
                    "success": True,
                    "total": tgbot_result["total"]
                }
            else:
                result = {
                    "success": False,
                    "message": "submit command, got last message: {}".format(last_message)
                }
        self.assert_successful(result)
