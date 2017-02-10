# licenser - adds missing copyright notice to all your source files

## Go
Adds your notice before the **package** statement using // as comment indicator.

## Usage
    licenser -d -go notice.txt

    -d : perform a dry run, see what files would be changed

    -go : files with extension .go will use this copyright notice

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