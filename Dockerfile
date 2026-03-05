# SPDX-License-Identifier: Apache-2.0

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

# build from node LTS Jod https://nodejs.org/en/about/releases/
FROM node:lts-jod@sha256:23c24e85395992be118734a39903e08c8f7d1abc73978c46b6bda90060091a49

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
