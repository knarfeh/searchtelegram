#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import datetime
import smtplib
import requests
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.header import Header

from maic.database.store import TestStore, TestCase
from sqlalchemy import create_engine, func
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship, sessionmaker
from api_test import configs


def get_case_num_db(case_name, status):
    host = os.getenv('DB_HOST', 'db')
    port = os.getenv('DB_PORT', '5432')
    username = os.getenv('DB_USER', 'knarfeh')
    password = os.getenv('DB_PASSWORD', 'no_password')
    dbname = os.getenv('DB_NAME', 'ste2e')

    Base = declarative_base()
    Session = sessionmaker()

    db_url = 'postgresql://{}:{}@{}:{}/{}?sslmode=disable'.format(
        username, password, host, port, dbname
    )

    engine = create_engine(db_url)
    Session.configure(bind=engine)
    session = Session()

    case_ended_time = session.query(TestCase.ended_at).\
        filter(TestCase.status == status,
               TestCase.name == case_name
               ).all()
    case_ended_time = case_ended_time[-1]
    timedelta = datetime.timedelta(days=1)
    one_day_before_ended = case_ended_time[0] - timedelta

    failed_case_num = session.query(func.count(TestCase.name)).\
        filter(TestCase.status == 'failed',
               TestCase.name == case_name,
               TestCase.ended_at <= case_ended_time,
               TestCase.started_at >= one_day_before_ended
               ).first()
    print('Today, {} case failed times : {}'.format(case_name, failed_case_num[0]))
    return failed_case_num[0]

def send_email(subject, body, recipients):
    sender = configs.EMAIL_CONFIG['sender']

    requests.post(
        "https://api.mailgun.net/v3/"+configs.MAILGUN_CONFIG['MAIL_DOMAIN_NAME']+'/messages',
        auth=("api", configs.MAILGUN_CONFIG["MAIL_API_KEY"]),
        data={
            "from": sender,
            "to": [recipients,],
            "subject": subject,
            "html": body
        }
    )
    return
