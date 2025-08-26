#!/usr/bin/env bash
# Copyright (c) 2025-2025 All rights reserved.

# The original source code is licensed under the Apache License 2.0.
# New features or modifications to this source code are licensed under the
# Affero General Public License 3.0 (AGPL 3.0).

# You may review the terms of both licenses in the LICENSE file.


set -xe

addlicense -f licenses.tpl -ignore "**/*.md" -ignore "**/*.rb" -ignore "vendor/**" -ignore "**/*.yml" -ignore "**/*.yaml" -ignore "**/*.sh" -ignore "**/*.sql" -ignore "**/*.html" ./**
