#!/bin/bash
#
# Copyright 2020 Aletheia Ware LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
set -x

(cd $GOPATH/src/github.com/AletheiaWareLLC/labfynego/ui/data/ && ./gen.sh)
go fmt $GOPATH/src/github.com/AletheiaWareLLC/{labfynego,labfynego/...}
go vet $GOPATH/src/github.com/AletheiaWareLLC/{labfynego,labfynego/...}
go test $GOPATH/src/github.com/AletheiaWareLLC/{labfynego,labfynego/...}
ANDROID_NDK_HOME=${ANDROID_HOME}/ndk-bundle/
(cd $GOPATH/src/github.com/AletheiaWareLLC/labfynego/cmd && fyne package -os android -appID com.aletheiaware.lab -icon $GOPATH/src/github.com/AletheiaWareLLC/labfynego/ui/data/logo.png -name LAB_unaligned)
#03-25 21:17:21.397 21848 21848 I letheiaware.la: Late-enabling -Xcheck:jni
#03-25 21:17:21.418 21848 21848 E letheiaware.la: Unknown bits set in runtime_flags: 0x8000
#03-25 21:17:21.448 21848 21848 W letheiaware.la: Can't mmap dex file /data/app/com.aletheiaware.lab-naC86cdO6jg6E9PTJnbWag==/base.apk!classes.dex directly; please zipalign to 4 bytes. Falling back to extracting file.
(cd $GOPATH/src/github.com/AletheiaWareLLC/labfynego/cmd && ${ANDROID_HOME}/build-tools/28.0.3/zipalign -f 4 LAB_unaligned.apk LAB.apk)
(cd $GOPATH/src/github.com/AletheiaWareLLC/labfynego/cmd && adb install -r -g LAB.apk)
#(cd $GOPATH/src/github.com/AletheiaWareLLC/labfynego/cmd && adb logcat com.aletheiaware.lab:V org.golang.app:V *:S | tee log)
(cd $GOPATH/src/github.com/AletheiaWareLLC/labfynego/cmd && adb logcat -c && adb logcat | tee android.log)
