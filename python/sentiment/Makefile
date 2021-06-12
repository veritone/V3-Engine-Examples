# Adjust this Makefile to suit your engine's needs.

# Name of this engine as provided by the Veritone Developer App under the build instructions
ENGINE_NAME = "YOUR-ENGINE-NAME"

#
# BUILDING targets
#
.PHONY: build clean deepclean install compile test

## Do a clean build of the docker image. 
build: clean install engine test docker

## Deletes intermediate build artifacts
clean: 
	rm -rf __pycache__ .pytest_cache

## Deletes all renewable artifacts. Restore the directory to as close to newly-cloned as possible
deepclean: ;

## Install dependencies
install: 
	pip3 install -r requirements.txt
	python3 -m nltk.downloader vader_lexicon
	python3 -m nltk.downloader punkt

## Compiles intermediate artifacts
compile: ;

## Run unit tests
test: install
	pytest -rP app_test.py

#
# ENGINE targets
#
.PHONY: engine docker

## Build the engine
engine: compile

## Build the docker image
docker: 
	docker build --tag=${ENGINE_NAME} .
