# Copyright 2022 Nexient LLC
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at

#   http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.    

REPO_MANIFESTS_URL ?= https://github.com/nexient-llc/module-manifests.git
REPO_BRANCH ?= main
REPO_MANIFEST ?= tf_modules.xml

# Set to true in a pipeline context
IS_PIPELINE ?= false

JOB_NAME ?= job
JOB_EMAIL ?= job@job.job

COMPONENTS_DIR = components
-include $(COMPONENTS_DIR)/Makefile

.PHONY: configure-git-hooks
configure-git-hooks:
	pre-commit install

ifeq ($(IS_PIPELINE),true)
.PHONY: git-config
git-config:
	$(call config,Bearer $(GIT_TOKEN))

define config
	@set -ex; \
	git config --global http.extraheader "AUTHORIZATION: $(1)"; \
	git config --global http.https://gerrit.googlesource.com/git-repo/.extraheader ''; \
	git config --global http.version HTTP/1.1; \
	git config --global user.name "$(JOB_NAME)"; \
	git config --global user.email "$(JOB_EMAIL)"; \
	git config --global color.ui false
endef

configure: git-config
endif

.PHONY: configure
configure: configure-git-hooks
	repo init \
		-u "$(REPO_MANIFESTS_URL)" \
		-b "$(REPO_BRANCH)" \
		-m "$(REPO_MANIFEST)"
	repo sync

# The first line finds and removes all the directories pulled in by repo
# The second line finds and removes all the broken symlinks from removing things
# https://stackoverflow.com/questions/42828021/removing-files-with-rm-using-find-and-xargs
.PHONY: clean
clean:
	repo list | awk '{ print $1; }' | cut -d '/' -f1 | uniq | xargs rm -rf
	find . -type l ! -exec test -e {} \; -print | xargs rm -rf
