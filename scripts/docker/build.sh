#!/usr/bin/env bash

set -e

# see https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds
export DOCKER_BUILDKIT=1

if [[ -z $1 ]] || [[ -z $2 ]]; then
    echo "Usage: ${0##*/} [name] [tag] [/subimage]" 1>&2
    exit 1
fi

NUMERIC='^[0-9]+$'
BUILD_DATE=$(/bin/date -u +%y%m%d)

echo "Building image 'hello/$1'";

if [[ $1 ]] && [[ -z $2 || $2 == "preview" ]]; then
    echo "Build Tags: preview"

    docker build \
      --no-cache \
      --build-arg BUILD_TAG=$BUILD_DATE \
      -t hello/$1:preview \
      -f docker/${1/-//}$3/Dockerfile .
elif [[ $2 =~ $NUMERIC ]]; then
    echo "Build Tags: $2, latest"

    fi

    docker build \
      --no-cache \
      --build-arg BUILD_TAG=$2 \
      -t hello/$1:latest \
      -t hello/$1:$2 \

elif [[ $2 == *"preview"* || $2 == *"unstable"* || $2 == *"test"* || $2 == *"local"* || $2 == *"develop"* ]]; then
    echo "Build Tags: $2"

    docker build $4\
      --no-cache \
      --build-arg BUILD_TAG=$BUILD_DATE \
      -t hello/$1:$2 \

else
    echo "Build Tags: $BUILD_DATE-$2, $2"

    docker build $4\
      --no-cache \
      --build-arg BUILD_TAG=$BUILD_DATE-$2 \
      -t hello/$1:$2 \
      -t hello/$1:$BUILD_DATE-$2 \

fi

echo "Done."
