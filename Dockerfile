# SPDX-License-Identifier: Apache-2.0

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

# build from node LTS Jod https://nodejs.org/en/about/releases/
FROM node:lts-jod@sha256:c3ef15af9be4505fde55589eadf42b4757a91e6b1b3be796bdec0f86560205e9

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
