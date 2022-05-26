/*
Copyright 2022 Opstree Solutions.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8selastic

const healthCheckScript = `#!/usr/bin/env bash -e
START_FILE=/tmp/.es_start_file
http () {
  local path="${1}"
  local args="${2}"
  set -- -XGET -s
  if [ "$args" != "" ]; then
    set -- "$@" $args
  fi
  if [ -n "${ELASTIC_USERNAME}" ] && [ -n "${ELASTIC_PASSWORD}" ]; then
    set -- "$@" -u "${ELASTIC_USERNAME}:${ELASTIC_PASSWORD}"
  fi
  curl --output /dev/null -k "$@" "${SCHEME}://127.0.0.1:9200${path}"
}
if [ -f "${START_FILE}" ]; then
  echo 'Elasticsearch is already running, lets check the node is healthy'
  HTTP_CODE=$(http "/" "-w %{http_code}")
  RC=$?
  if [[ ${RC} -ne 0 ]]; then
    echo "curl --output /dev/null -k -XGET -s -w '%{http_code}' \${BASIC_AUTH} ${SCHEME}://127.0.0.1:9200/ failed with RC ${RC}"
    exit ${RC}
  fi
  # ready if HTTP code 200, 503 is tolerable if ES version is 6.x
  if [[ ${HTTP_CODE} == "200" ]]; then
    exit 0
  elif [[ ${HTTP_CODE} == "503" && "7" == "6" ]]; then
    exit 0
  else
    echo "curl --output /dev/null -k -XGET -s -w '%{http_code}' \${BASIC_AUTH} ${SCHEME}://127.0.0.1:9200/ failed with HTTP code ${HTTP_CODE}"
    exit 1
  fi
else
  echo 'Waiting for elasticsearch cluster to become ready (request params: "wait_for_status=green&timeout=1s" )'
  if http "/_cluster/health?wait_for_status=green&timeout=1s" "--fail" ; then
    touch ${START_FILE}
    exit 0
  else
    echo 'Cluster is not yet ready (request params: "wait_for_status=green&timeout=1s" )'
    exit 1
  fi
fi
`
