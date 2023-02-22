run:
	rm -rf TEST_PROJECT
	go run main.go \
		--out=./gen \
		--no-clobber=true \
		--log-level=debug \
		new .scaffolds/cli \
		"Project=TEST_PROJECT" \
		"Description=TEST_PROJECT" \
		"License=MIT" \
		"Colors=#000000" \
		"Use Github Actions=true"