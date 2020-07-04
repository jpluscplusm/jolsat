test: test-binaries test-no-args test-one-arg-d
	@echo "** All tests passed"

test-binaries: jolsat
	echo | ./jolsat
	echo | cat

test-no-args: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -f1- | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat | md5sum)" ]

test-one-arg-d: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f1- | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' | md5sum)" ]
