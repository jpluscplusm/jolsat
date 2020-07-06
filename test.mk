test: test-binaries test-no-args test-one-arg-d test-two-args-d-f-single-field test-two-args-d-f-range test-two-args-d-f-range-and-field test-two-args-d-f-multiple-ranges test-two-args-d-f-range-and-multiple-fields
	@echo "** All tests passed"

test-binaries: jolsat
	@echo | ./jolsat >/dev/null
	@echo | cat >/dev/null

test-no-args: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -f 1- | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat | md5sum)" ]

test-one-arg-d: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1- | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' | md5sum)" ]

test-two-args-d-f-single-field: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 3 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 3 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 5 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 5 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 6 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 6 | md5sum)" ]

test-two-args-d-f-range: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-1 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-1 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-2 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-2 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-5 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-5 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 2-4 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 2-4 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 2-5 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 2-5 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 2-6 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 2-6 | md5sum)" ]

test-two-args-d-f-range-and-field: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-3,5 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-3,5 | md5sum)" ]

test-two-args-d-f-multiple-ranges: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-2,4-5 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-2,4-5 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-2,4-6 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-2,4-6 | md5sum)" ]

test-two-args-d-f-range-and-multiple-fields: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-2,3,4,5 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-2,3,4,5 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1-2,3,4,5,6 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1-2,3,4,5,6 | md5sum)" ]
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1,3-4,6 | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f 1,3-4,6 | md5sum)" ]
