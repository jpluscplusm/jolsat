test: jolsat
	echo | ./jolsat
	echo | cat
	[ "$(shell seq 1 20 | xargs -n5 | cat | md5sum)" = "$(shell seq 1 20 | xargs -n5 | ./jolsat | md5sum)" ]

jolsat:	jolsat.go
	go build
