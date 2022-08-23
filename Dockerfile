# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

# build from node LTS Gallium https://nodejs.org/en/about/releases/
FROM node:lts-gallium@sha256:d8d3181ca9840f6667ae8694c35511af806e31c45f9c4fa5f80328b0b2c1dc44

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
