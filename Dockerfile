# SPDX-License-Identifier: Apache-2.0

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

# build from node LTS Jod https://nodejs.org/en/about/releases/
FROM node:lts-jod@sha256:263acfc82e6c90dc5f80a937868a7ede7f205f1e0a1eecd94b254e3497829dbb

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
