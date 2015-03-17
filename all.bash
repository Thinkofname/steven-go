#!/usr/bin/env bash
# Copyright 2014 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Script to build and launch the app on an android device.

set -e

# ./make.bash
docker run -v $GOPATH/src:/src golang/mobile /bin/bash -c 'cd /src/github.com/thinkofdeath/steven && ./make.bash'

adb install -r bin/Steven-debug.apk

adb shell am start -a android.intent.action.MAIN \
	-n uk.co.thinkofdeath.steven/android.app.NativeActivity
