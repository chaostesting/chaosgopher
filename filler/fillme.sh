#!/bin/sh

# REF: http://www.commandlinefu.com/commands/view/4249/fill-up-disk-space-for-testing
echo "FILLING UP SPACE..."
yes > /filler.dat || true
echo "FILLED UP, SLEEPING..."
sleep 86400
