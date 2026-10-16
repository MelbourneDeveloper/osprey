# Build from the repository root. See docs/contributing/repository-map.md.
.DEFAULT_GOAL := build

include scripts/make/config.mk
include scripts/make/targets.mk
include scripts/make/runtime.mk
include scripts/make/tests.mk
include scripts/make/apps.mk
include scripts/make/extension.mk

include scripts/ios.mk
include scripts/android.mk
include scripts/bank-mobile.mk
