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

package k8sgo

const keyStoreCommand = `set -euo pipefail
elasticsearch-keystore create
for i in /tmp/keystoreSecrets/*/*; do
  key=$(basename $i)
  echo "Adding file $i to keystore key $key"
  elasticsearch-keystore add-file "$key" "$i"
done
# Add the bootstrap password since otherwise the Elasticsearch entrypoint tries to do this on startup
if [ ! -z ${ELASTIC_PASSWORD+x} ]; then
  echo 'Adding env $ELASTIC_PASSWORD to keystore as key bootstrap.password'
  echo "$ELASTIC_PASSWORD" | elasticsearch-keystore add -x bootstrap.password
fi
cp -a /usr/share/elasticsearch/config/elasticsearch.keystore /tmp/keystore/`
