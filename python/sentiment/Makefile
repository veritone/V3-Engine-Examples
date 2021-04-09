# Template Makefile
# Adjust this Makefile to suit your engine's needs, but this template includes

# Name of this engine as provided by the Veritone Developer App under the build instructions
VDA_NAME = "sentiment-sample-chunk-v-3"

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
.PHONY: engine docker test-image publish deploy deploy-%

## Build the engine
engine: compile

## Build the docker image
docker: 
	docker build --tag=${VDA_NAME} .

# start the engine in the test framework
up-test:
	docker run --rm -e "VERITONE_TESTMODE=true" -p 9090:9090 -p 8080:8080 $(VDA_NAME)

# rebuild and start the engine in the test framework
up-retest: docker up-test
