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
go run -tags mobile github.com/AletheiaWareLLC/labfynego/cmd $@