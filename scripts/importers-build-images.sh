#!/usr/bin/env bash

(cd user-importer && docker build -t textsearch/userimport .)
(cd relevance-importer && docker build -t textsearch/relevanceimport .)