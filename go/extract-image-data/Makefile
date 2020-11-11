current_dir = $(shell pwd)

# Do a clean build of the docker image
build: clean engine docker

# Convenience target to shut down current engine, recompile, reimage, and run it again
run: down engine docker up

# Convenience target to shut down current test engine, recompile, reimage, and launch test engine again
run-testmode: down engine docker up-testmode

# Convenience target to shut down current engine, recompile, reimage, and run the engine with the file service
run-fs: down engine docker up-fs

engine:
	GOOS=linux GOARCH=amd64 go build -o ./dist/exif-extraction-engine

clean:
	rm -rf ./dist

docker:
	docker build -f Dockerfile -t exif-extraction-engine .

up:
	docker run -t exif-extraction-engine --name exif-extraction-engine

down:
	docker rm -f exif-extraction-engine 2>/dev/null || true

up-testmode:
	docker run --rm -e "VERITONE_TESTMODE=true" -p 9090:9090 -p 8080:8080 --name exif-extraction-engine -t exif-extraction-engine

up-fs:
	docker run \
		-v $(current_dir)/testdata/fs-engine-in:/files/in \
		-v $(current_dir)/testdata/fs-engine-completed:/files/out/completed \
		-v $(current_dir)/testdata/fs-engine-errors:/files/out/errors \
		-v $(current_dir)/testdata/fs-engine-results:/files/out/results \
		-e "VERITONE_SELFDRIVING=true" \
		-e "VERITONE_SELFDRIVING_INPUTPATTERN=*.jpg" \
		-e "VERITONE_SELFDRIVING_OUTPUT_DIR_PATTERN=yyyy/mm/dd" \
		-e "VERITONE_SELFDRIVING_MINIMUM_MODIFIED_DURATION=1s" \
		-t exif-extraction-engine --name exif-extraction-engine
