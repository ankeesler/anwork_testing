#
# Makefile
#
# Andrew Keesler
#
# November 25, 2017
# Caribou Coffee, Park Road Shopping Center, with my sister sitting across the table
#
# This is the main Makefile for the anwork_testing project. It provides shortcut targets for doing
# stuff like
#   running all the tests    (make test)
#   running a single test    (make test-XXX)
#   generating code coverage (make coverage)
#   generating code profile  (make profile-XXX)
#   cleaning test output     (make clean)
#

PACKAGES=core v1

.PHONY: test
test:
	GOPATH="$(PWD)" go test $(PACKAGES)

test-%:
	GOPATH="$(PWD)" go test $(patsubst test-%,%,$@)

.PHONY: coverage
coverage:
	GOPATH="$(PWD)" go test core -coverprofile=cover.out
	go tool cover -html=cover.out

.PHONY:
profile-%:
	GOPATH="$(PWD)" go test core -cpuprofile=cpu.out
	go tool pprof -web cpu.out

.PHONY: clean
clean:
	rm -rf .anwork-*
	rm -f *.out
	rm -f $(patsubst %,%.test,$(PACKAGES))