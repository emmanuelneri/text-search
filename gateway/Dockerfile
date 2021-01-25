FROM kong:2.3.0
USER root
RUN apk update && apk add git unzip luarocks
RUN luarocks install kong-oidc
USER kong