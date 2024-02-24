VERSION 0.7

FROM tochemey/docker-go:1.21.0-1.0.0

RUN apk --no-cache add git ca-certificates gcc musl-dev libc-dev binutils-gold

protogen:
    # copy the proto files to generate
    COPY --dir protos/ ./
    COPY buf.work.yaml buf.gen.yaml ./

    # generate the pbs
    RUN buf generate \
            --template buf.gen.yaml \
            --path protos/groupcache \
            --path protos/example \
            --path protos/test

    # save artifact to
    SAVE ARTIFACT gen/groupcache AS LOCAL groupcachepb
    SAVE ARTIFACT gen/test AS LOCAL testpb
    SAVE ARTIFACT gen/example AS LOCAL example/pb
