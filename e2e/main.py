#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = "knarfeh@outlook.com"

import logging
import json

from time import mktime, strptime
from api_test import ApiTest
from api_test import config
from utils import get_case_num_db, send_email

LOGGER = logging.getLogger()


def main():
    """
    Main function to start testing
    """
    print("hello")
    test = ApiTest(name="api", environment="dev", workers=4, timeout=60*60, log=True)
    test.set_up()

    LOGGER.info("Starting test...")
    test.start()

    LOGGER.info("End of test...")
    result_data = test.end()

    LOGGER.info("Tear down")
    test.tear_down()
    LOGGER.info("Done! \n")
    print('=================RESULT===================')
    print('Result of test.end: ' + json.dumps(result_data))
    result_status = result_data['status']
    print('Result status: ' + result_status)


if __name__ == "__main__":
    main()
