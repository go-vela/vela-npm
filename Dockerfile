# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

###########################################################################
##    docker build --no-cache -t vela-npm:local .    ##
###########################################################################

FROM node:lts-gallium@sha256:2e1b4542d4a06e0e0442dc38af1f4828760aecc9db2b95e7df87f573640d98cd

COPY release/vela-npm /bin/vela-npm

ENTRYPOINT [ "/bin/vela-npm" ]
