# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

# build from node LTS Gallium https://nodejs.org/en/about/releases/
FROM node:lts-gallium@sha256:1ed1e17ccabb09038cfb8a965337ebcda51ef9e9d32082164c502d44d9731a02

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
