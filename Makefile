# Copyright 2019 Philippe Martin
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

default:
	@echo "commands: clean, docbook, pdf"

clean:
	rm -rf build

docbook: clean build/index.xml

build/index.xml: $(wildcard *.go **/*.go) generators/v1_17/toc.yaml
	mkdir -p build
	LANG= go run main.go --kubernetes-version v1_17 > build/index.xml

FORMAT ?= USletter
pdf: build/index.xml
	(cd build && \
	mkdir -p pdf-$(FORMAT) && \
	cd pdf-$(FORMAT) && \
	xsltproc --stringparam fop1.extensions 1 --stringparam paper.type $(FORMAT) -o index-$(FORMAT).fo ../../xsl/api.xsl ../index.xml && \
	fop -pdf index-$(FORMAT).pdf -fo index-$(FORMAT).fo && \
	rm  index-$(FORMAT).fo)

pdf-6x9in: build/index.xml
	(cd build && \
	mkdir -p pdf && \
	cd pdf && \
	xsltproc --stringparam fop1.extensions 1 -o index.fo ../../xsl/api-6x9in.xsl ../index.xml && \
	fop -pdf index.pdf -fo index.fo && \
	rm  index.fo)

test:
	@echo $(FORMAT)
