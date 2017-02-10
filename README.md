# licenser - adds copyright notice to all your source files

[![Build Status](https://travis-ci.org/emicklei/licenser.png)](https://travis-ci.org/emicklei/licenser)

    go get github.com/emicklei/licenser

## Help

    Usage: licenser [flags] [path...]

    -d	dry run, see what would change
    -e string
            file extension for which the copyright notice must be added (default ".go")
    -f string
            source token that indicates where the actual source will start
    -t string
            source token that indicates where the actual source will start (default "package")

## Example

    licenser -d -e ".go" -f LICENSE -t "package" .

## Example Apache v2

    Copyright 2017 Me, myself and I

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.