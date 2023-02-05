run:
	rm -rf TEST_PROJECT
	go run main.go \
		--cwd=./gen \
		--log-level=debug \
		new .scaffolds/cli