###############################################################################
# Copyright 2017 Samsung Electronics All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
###############################################################################
# Docker image for "Service Deployment Agent"
FROM alpine:3.5
MAINTAINER Heejung Kim <hjnari.kim@samsung.com>

# environment variables
ENV APP_DIR=/ServiceDeploymentAgentManager
ENV APP_PORT=48099

# repositories for mongodb
RUN echo 'http://dl-3.alpinelinux.org/alpine/v3.6/community' >> /etc/apk/repositories
RUN echo 'http://dl-3.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories

# install docker-compose
RUN apk update
RUN apk add curl

# install MongoDB
RUN apk add --no-cache --update mongodb && \
    rm -rf /var/cache/apk/*

#copy files
COPY main $APP_DIR/
COPY run.sh $APP_DIR

#expose notifications port
EXPOSE $APP_PORT

#set the working directory
WORKDIR $APP_DIR

#make mogodb volume
RUN mkdir /data
RUN mkdir /data/db
VOLUME /data/db

#kick off the agent container
CMD ["sh", "run.sh"]
