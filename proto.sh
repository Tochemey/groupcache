#! /bin/sh

# Make sure the script fails fast.
set -e
set -u
set -x

PROTO_DIR=groupcachepb
PROTO_TEST_DIR=testpb
NATS_PROTO_DIR=discovery/nats

protoc -I=$PROTO_DIR \
    --go_out=$PROTO_DIR \
    $PROTO_DIR/groupcache.proto

protoc -I=$PROTO_TEST_DIR \
    --go_out=$PROTO_TEST_DIR \
    $PROTO_TEST_DIR/test.proto

protoc -I=. \
   --go_out=. \
    example.proto

protoc -I=$NATS_PROTO_DIR \
    --go_out=$NATS_PROTO_DIR \
    $NATS_PROTO_DIR/nats.proto
