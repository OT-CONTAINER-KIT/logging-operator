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

package k8skibana

const healthCheckScript = `#!/usr/bin/env bash -e
http () {
	local path="${1}"
	set -- -XGET -s --fail
	if [ -n "${ELASTIC_USERNAME}" ] && [ -n "${ELASTIC_PASSWORD}" ]; then
	  set -- "$@" -u "${ELASTIC_USERNAME}:${ELASTIC_PASSWORD}"
	fi
	STATUS=$(curl --output /dev/null --write-out "%{http_code}" -k "$@" "http://localhost:5601${path}")
	if [[ "${STATUS}" -eq 200 ]]; then
	  exit 0
	fi
	echo "Error: Got HTTP code ${STATUS} but expected a 200"
	exit 1
}
http "/app/kibana"
`
