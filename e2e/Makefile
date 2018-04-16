OS = Linux

VERSION = 0.0.1

CURDIR = $(shell pwd)
SOURCEDIR = $(CURDIR)

ECHO = echo
RM = rm -rf
MKDIR = mkdir
FLAKE8 = flake8
PIP_INSTALL = pip install

.PHONY: all setup build test help

all: build test

setup:
	$(PIP_INSTALL) -r requirements/dev.txt

build:
	$(FLAKE8) $(SOURCEDIR) --show-source --show-pep8 --statistics --count

test:
	python -m unittest discover .

coverage:
	nosetests --cover-tests --cover-inclusive --cover-package=alauda

help:
	@$(ECHO) "Targets:"
	@$(ECHO) "all      - Run build and test"
	@$(ECHO) "setup    - Set up prerequisition for build"
	@$(ECHO) "build    - Run flake8"
	@$(ECHO) "test     - Run tests"
	@$(ECHO) "coverage - Work in progress"
