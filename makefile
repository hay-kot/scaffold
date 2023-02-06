run:
	rm -rf TEST_PROJECT
	go run main.go \
		--cwd=./gen \
		--log-level=debug \
		--var "Project=TEST_PROJECT" \
		--var "Description=TEST_PROJECT" \
		--var "License=MIT" \
		--var "Colors=#000000" \
		--var "Use Github Actions=true" \
		new .scaffolds/cli