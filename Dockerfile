# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

FROM node@sha256:82f9e078898dce32c7bf3232049715f1b8fbf0d62d5f3091bca20fcaede50bf0

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
