#!/bin/bash
set -e

#-d         no daemon, run as cli
#-U root    run in root context
#-l stdout  print out to stdout

maxscale -d -U root -l file
