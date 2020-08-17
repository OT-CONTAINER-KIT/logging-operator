package deployment

const readinessScript = `#!/usr/bin/env bash -e
http () {
	local path="${1}"
	set -- -XGET -s --fail

	if [ -n "${ELASTICSEARCH_USERNAME}" ] && [ -n "${ELASTICSEARCH_PASSWORD}" ]; then
	  set -- "$@" -u "${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}"
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
