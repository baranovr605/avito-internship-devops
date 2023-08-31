#!/bin/bash

# Generate some test certificates which are used by the regression test suite:
#
#   certs/ca.{crt,key}          Self signed CA certificate.
#   certs/redis.{crt,key}       A certificate with no key usage/policy restrictions.
#   certs/redis.dh              DH Params file.

CERTS_DIR="./certs"

generate_cert() {
    local name=$1
    local cn="$2"
    local opts="$3"

    local keyfile=${CERTS_DIR}/${name}.key
    local certfile=${CERTS_DIR}/${name}.crt

    [ -f $keyfile ] || openssl genrsa -out $keyfile 2048
    openssl req \
        -new -sha256 \
        -subj "/O=Redis Test/CN=$cn" \
        -key $keyfile | \
        openssl x509 \
            -req -sha256 \
            -CA ${CERTS_DIR}/ca.crt \
            -CAkey ${CERTS_DIR}/ca.key \
            -CAserial ${CERTS_DIR}/ca.txt \
            -CAcreateserial \
            -days 365 \
            $opts \
            -out $certfile
}

mkdir mkdir -p ${CERTS_DIR}
[ -f ${CERTS_DIR}/ca.key ] || openssl genrsa -out ${CERTS_DIR}/ca.key 4096
openssl req \
    -x509 -new -nodes -sha256 \
    -key ${CERTS_DIR}/ca.key \
    -days 3650 \
    -subj '/O=Redis Test/CN=Certificate Authority' \
    -out ${CERTS_DIR}/ca.crt

generate_cert redis "Generic-cert"

[ -f ${$CERTS_DIR}/redis.dh ] || openssl dhparam -out ${CERTS_DIR}/redis.dh 2048
