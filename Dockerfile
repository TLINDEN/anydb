# 
# Copyright © 2024 Thomas von Dein
# 
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.
# 
FROM golang:1.24-alpine as builder

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git make

RUN git --version

WORKDIR /work

COPY go.mod .
COPY . .
RUN go mod download
RUN make

FROM alpine:latest
LABEL maintainer="Thomas von Dein <git@daemon.de>"

WORKDIR /app
COPY --from=builder /work/anydb /app/anydb

ENV LANG C.UTF-8
USER 1001:1001

ENTRYPOINT ["/app/anydb"]
CMD ["-h"]
