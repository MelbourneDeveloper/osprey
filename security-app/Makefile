OSPREY_COMPILER ?= osprey
export OSPREY_COMPILER

.PHONY: build test ci ios
build:
	./build.sh

test: build
	"$(OSPREY_COMPILER)" engine/tests --run
	swift test

ios:
	xcodebuild -scheme Vigil -destination 'generic/platform=iOS' -derivedDataPath .build/ios CODE_SIGNING_ALLOWED=NO build

ci: test ios
