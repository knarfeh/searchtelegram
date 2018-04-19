#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = "knarfeh@outlook.com"

import logging
import json

from time import mktime, strptime, sleep
from api_test import ApiTest
from api_test import configs
from api_test.api_tgbot import get_tgbot_response
from utils import get_case_num_db, send_email, get_tgbot_payload

LOGGER = logging.getLogger()


def main():
    """
    Main function to start testing
    """
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

    sleep(10)
    payload = get_tgbot_payload("/echo e2e tests result: {}".format(result_status))
    get_tgbot_response(payload, configs.get_header())

    e2e_start_time = mktime(strptime(result_data['started'][0:-6], "%Y-%m-%dT%H:%M:%S.%f"))
    print('e2e_start_time: {}'.format(e2e_start_time))
    e2e_end_time = mktime(strptime(result_data['ended'][0:-6], "%Y-%m-%dT%H:%M:%S.%f"))
    print('e2e_end_time: {}'.format(e2e_end_time))
    total_time = e2e_end_time - e2e_start_time
    print('total_time: {}'.format(total_time))

    case_name = []
    case_details = []
    case_nums_failed_db = []
    cases_time = []
    cases = result_data['cases']
    cases_flags = []
    total_case_num = len(cases)
    failed_cases_num = 0

    for i in range(len(cases)) :
        caseflag = cases[i]['status']
        print("Running cases name: " + cases[i]['name'])
        case_start_time = mktime(strptime(cases[i]['started'][0:-6], "%Y-%m-%dT%H:%M:%S.%f"))
        print('Case start time: {}'.format(case_start_time))
        case_end_time = mktime(strptime(cases[i]['ended'][0:-6], "%Y-%m-%dT%H:%M:%S.%f"))
        print('Case end time : {}'.format(case_end_time))
        case_total_time = case_end_time - case_start_time
        print('Case total time : {}'.format(case_total_time))

        case_details.append(json.dumps(cases[i]['details']))
        case_name.append(cases[i]['name'])
        cases_flags.append(caseflag)
        cases_time.append(case_total_time)

        if caseflag == "failed" :
            failed_cases_num = failed_cases_num + 1
            print(cases[i])

        case_nums_failed_db.append(get_case_num_db(cases[i]['name'], caseflag))

    print('Total failed cases number : {}'.format(failed_cases_num))

    html = "Run {} cases, Pass {} cases, Fail {} cases, ".format(total_case_num, total_case_num - failed_cases_num, failed_cases_num) + "\n" + \
           "e2e starts at {}, ".format(result_data['started'][0:-13]) + "\n" + \
           "ends at {}, ".format(result_data['ended'][0:-13]) + "\n" + \
           "Total time {}s ".format(total_time) + "\n"

    html = html + '<table border="1">\n' \
                  '<thead>\n' \
                  '<tr>\n' \
                  '<td>Case name</td>\n' \
                  '<td>Case flag</td>\n' \
                  '<td>Case details</td>\n' \
                  '<td>Case time</td>\n' \
                  '<td>Case failed times (today)</td>' \
                  '\n</tr>\n' \
                  '</thead>\n' \
                  '<tbody>'

    for i in range(len(case_name)) :
        if cases_flags[i] == 'failed':
            html = html + '\n<tr>\n' \
                          '<td>{}</td>'.format(case_name[i]) + \
                    '\n<td style="color:red;">{}</td>'.format(cases_flags[i]) + \
                   '\n<td style="color:red;">{}</td>'.format(case_details[i]) + \
                   '\n<td style="color:red;">{}s</td>'.format(cases_time[i]) + \
                   '\n<td style="color:red;">{}</td>'.format(case_nums_failed_db[i]) + '\n</tr>'
        else :
            html = html + '\n<tr>\n<td>{}</td>'.format(case_name[i]) + \
                   '\n<td style="color:green;">{}</td>'.format(cases_flags[i]) + \
                   '\n<td style="color:green;">{}</td>'.format(case_details[i]) + \
                   '\n<td style="color:green;">{}s</td>'.format(cases_time[i]) + \
                   '\n<td style="color:green;">{}</td>'.format(case_nums_failed_db[i]) + '\n</tr>'

    html = html + '\n</tbody>\n</table>'

    # =========send email=============
    # send_email(
    #     "[{}] ({}) End-to-End Test".format(result_status.upper(), configs.ENV),
    #     html,
    #     configs.EMAIL_CONFIG['recipient'])

    # print case results in console
    print("--------------------- Case results -------------------")
    print(html)
    print("--------------------- Case results -------------------")

    f = open("/usr/share/nginx/html/index.html", 'w+')
    f.write(html)
    f.close()

    payload = get_tgbot_payload("/echo check out: https://e2e.searchtelegram.com")
    get_tgbot_response(payload, configs.get_header())


if __name__ == "__main__":
    main()
