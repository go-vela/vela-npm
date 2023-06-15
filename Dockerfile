# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

# build from node LTS Gallium https://nodejs.org/en/about/releases/
FROM node:lts-gallium@sha256:676b6c3c77f7d4324d1f8dceff33e3c8b08d9089016ab59c0657852aa95f9eb7

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
