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

fyne bundle -name Logo -package data lab.svg > icon.go
fyne bundle -append -name LogoUnmasked -package data lab-unmasked.svg >> icon.go
#fyne bundle -append -name FileNewIcon -package data file-new.svg >> icon.go
#fyne bundle -append -name CloudSaveIcon -package data cloud-save.svg >> icon.go
#fyne bundle -append -name PencilIcon -package data pencil.svg >> icon.go
#fyne bundle -append -name EraserIcon -package data eraser.svg >> icon.go
