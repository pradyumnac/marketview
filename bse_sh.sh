#!/usr/bin/env bash

marketview -h $1|column -t -s',' -c 40|fzf
