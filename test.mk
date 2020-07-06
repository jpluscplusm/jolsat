test: test-binaries test-no-args test-one-arg-d test-cut-compat test-cut-differences
	@echo "** All tests passed"

test-binaries: jolsat
	@echo | ./jolsat >/dev/null
	@echo | cat >/dev/null

test-no-args: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -f 1- | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat | md5sum)" ]

test-one-arg-d: jolsat
	[ "$(shell seq 1 20 | xargs -n5 | cut -d ' ' -f 1- | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat -d ' ' | md5sum)" ]

CUT_POSITIVE_COMPAT = 1 3 5 6 1-1 1-2 1-5 2-4 2-5 2-6 1-3,5 1-2,4-5 1-2,4-6 1-2,3,4,5 1-2,3,4,5,6 1,3-4,6
test-cut-compat: jolsat
	for FIELD in $(CUT_POSITIVE_COMPAT); do [ "$$(seq 1 20 | xargs -n5 | cut -d ' ' -f $${FIELD} | md5sum)" = "$$(seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f $${FIELD} | md5sum)" ] ; done

CUT_NEGATIVE_COMPAT = 1,1 2,1 1-2,3,3
test-cut-differences: jolsat
	for FIELD in $(CUT_NEGATIVE_COMPAT); do [ "$$(seq 1 20 | xargs -n5 | cut -d ' ' -f $${FIELD} | md5sum)" != "$$(seq 1 20 | xargs -n5 | ./jolsat -d ' ' -f $${FIELD} | md5sum)" ] ; done
