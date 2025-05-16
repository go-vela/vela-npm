# SPDX-License-Identifier: Apache-2.0

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

# build from node LTS Jod https://nodejs.org/en/about/releases/
FROM node:lts-jod@sha256:e558507eb799e3a76fcdaaee5e48dce1a00aebc85892128a9fca59f63bd49511

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
