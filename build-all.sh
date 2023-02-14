#!/bin/bash
# Copyright 2022 ok-john <github.com/ok-john>
# Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
# 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
# 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
#THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

VERSION="$(git tag|tail -n 1)"
echo "build with version tag: $VERSION"
readonly __HOST_ARCH=${1:-"arm64"}  # change this for your machine.
readonly __HOST_GOOSE=${2:-"darwin"} # change this for your machine.
readonly __BIN_NAME=${3:-"oai"}
readonly __EXEC_DIR=$(dirname "$(realpath $0)") && cd $__EXEC_DIR   

__darwin=( arm64 amd64 )
__linux=( amd64 arm arm64 mips mips64 mips64le mipsle ppc64 ppc64le riscv64 s390x 386 )
__freebsd=( amd64 arm arm64 386 )
__windows=( amd64 arm arm64 386 ) 

function build
{   
    local _goarch=${1:-"None"} && if [[ $_goarch == "None" ]]; then exit 1; fi
    local _goose="${2:-"None"}" && if [[ $_goose == "None" ]]; then exit 1; fi
    local _goarm=${3:-""}
    local _out="build/$_goose-$_goarch$_goarm" && mkdir -p $_out
    _out="$_out/$__BIN_NAME"
    if [ "$_goarch" == "arm" ] && [ "$_goarm" == "" ]; then
	    build $_goarch $_goose 5 && build $_goarch $_goose 6 && build $_goarch $_goose 7
    else
        
        if [[ $_goarch == mips* ]]; then
            #At present GOMIPS64 based binaries are not generated through this script, more details about GOMIPS environment variables in https://go.dev/doc/asm#mips .
            echo $_out-softfloat
            GOARM=$_goarm GOMIPS=softfloat GOARCH=$_goarch GOOS=$_goose GOHOSTARCH=$__HOST_ARCH CGO_ENABLED=0 go build -ldflags="-X 'main.version=$VERSION'" -o $_out-softfloat
            echo $_out
            GOARM=$_goarm GOARCH=$_goarch GOOS=$_goose GOHOSTARCH=$__HOST_ARCH CGO_ENABLED=0 go build -ldflags="-X 'main.version=$VERSION'" -o $_out
        else
            echo $_out
            GOARM=$_goarm GOARCH=$_goarch GOOS=$_goose GOHOSTARCH=$__HOST_ARCH CGO_ENABLED=0 go build -ldflags="-X 'main.version=$VERSION'" -o $_out
        fi
    fi
    sha256sum $_out > "$_out.sha256sum"
}

for arch in ${__linux[*]}; do build "$arch" "linux"; done

for arch in ${__freebsd[*]}; do build "$arch" "freebsd"; done

for arch in ${__darwin[*]}; do build "$arch" "darwin"; done

for arch in ${__windows[*]}; do build "$arch" "windows"; done
