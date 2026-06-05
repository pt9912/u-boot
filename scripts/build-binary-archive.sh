#!/bin/sh
set -eu

: "${UBOOT_VERSION:?UBOOT_VERSION is required}"
: "${PLATFORMS:?PLATFORMS is required}"

out_dir="${OUT_DIR:-/out}"
rm -rf "$out_dir"
mkdir -p "$out_dir"

for platform in $PLATFORMS; do
    os="${platform%/*}"
    arch="${platform#*/}"
    if [ "$os" = "$platform" ] || [ -z "$os" ] || [ -z "$arch" ]; then
        echo "invalid platform ${platform}; expected os/arch" >&2
        exit 2
    fi

    ext=""
    if [ "$os" = "windows" ]; then
        ext=".exe"
    fi

    out="$out_dir/u-boot-$os-$arch$ext"
    echo "==> $out (UBOOT_VERSION=$UBOOT_VERSION)" >&2
    GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 \
        go build \
        -ldflags="-s -w -X main.version=$UBOOT_VERSION" \
        -o "$out" \
        ./cmd/uboot
done

tar -C "$out_dir" -cf - .
