#!/bin/sh
# zig cc wrapper used for the autoconf/make build of jemalloc, a C dependency
# of rumdl, on the 32 bits Arm targets.
#
# cargo-zigbuild passes the C compiler the flags rustc itself receives, and zig
# does not understand two of them on arm/armv7:
# - --target=<rust triple>, for example --target=armv7-unknown-linux-musleabihf,
#   which is not a valid zig target query (zig wants arm-linux-musleabihf);
# - -march=armv7-a -mfpu=vfpv3-d16 -mfloat-abi=hard, because zig selects the Arm
#   architecture, FPU and float ABI with -mcpu=<CPU>+<feature>... instead, and
#   fails with 'unsupported argument generic+v7a+vfp3-d32+thumb2-neon to option
#   -mcpu=' when it deduces that string from the Rust flags.
# They are replaced by the zig target triple and -mcpu of the current build,
# read from the environment, which keep the exact same architecture baseline as
# rustc uses ('+v7a' and '+v6' are zig feature names, and plain -mcpu=generic
# would mean ARMv4, too old to build musl).
#
# Nothing here is specific to rumdl: the same wrapper works for any Rust program
# of this repository needing it, only the two environment variables change. See
# the notes on Rust programs in the README.
args=
for arg in "$@"; do
	case "${arg}" in
	--target=* | -march=* | -mcpu=* | -mfpu=* | -mfloat-abi=*) continue ;;
	*) args="${args}${args:+ }${arg}" ;;
	esac
done
# shellcheck disable=SC2086
exec zig cc -target "${ZIG_CC_WRAPPER_TARGET}" \
	${ZIG_CC_WRAPPER_MCPU:+-mcpu=${ZIG_CC_WRAPPER_MCPU}} ${args}
