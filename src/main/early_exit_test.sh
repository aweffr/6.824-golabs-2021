#!/usr/bin/env bash

#
# basic map-reduce test
#

#RACE=

# comment this to run the tests without the Go race detector.
RACE=-race

# run the test in a fresh sub-directory.
rm -rf mr-tmp
mkdir mr-tmp || exit 1
cd mr-tmp || exit 1
rm -f mr-*

# make sure software is freshly built.
(cd ../../mrapps && go build $RACE -buildmode=plugin wc.go) || exit 1
(cd ../../mrapps && go build $RACE -buildmode=plugin early_exit.go) || exit 1
(cd .. && go build $RACE mrcoordinator.go) || exit 1
(cd .. && go build $RACE mrworker.go) || exit 1

failed_any=0
#########################################################
# test whether any worker or coordinator exits before the
# task has completed (i.e., all output files have been finalized)
echo '***' Starting early exit test.

timeout -k 2s 180s ../mrcoordinator ../pg*txt &

# give the coordinator time to create the sockets.
sleep 1

# start multiple workers.
timeout -k 2s 180s ../mrworker ../../mrapps/early_exit.so &
timeout -k 2s 180s ../mrworker ../../mrapps/early_exit.so &
timeout -k 2s 180s ../mrworker ../../mrapps/early_exit.so &

# wait for any of the coord or workers to exit
# `jobs` ensures that any completed old processes from other tests
# are not waited upon
jobs &> /dev/null
wait -n

# a process has exited. this means that the output should be finalized
# otherwise, either a worker or the coordinator exited early
sort mr-out* | grep . > mr-wc-all-initial

# wait for remaining workers and coordinator to exit.
wait

# compare initial and final outputs
sort mr-out* | grep . > mr-wc-all-final
if cmp mr-wc-all-final mr-wc-all-initial
then
  echo '---' early exit test: PASS
else
  echo '---' output changed after first worker exited
  echo '---' early exit test: FAIL
  failed_any=1
fi
