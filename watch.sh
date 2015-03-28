#!/bin/bash -e
#
# Runs in the container host (i.e. CoreOS)
#
# ETCD_DEPLOY_KEY=/containers/[quay repository url]/latest/build

etcdctl exec-watch $ETCD_DEPLOY_KEY -- sh -c "sh run.sh"
