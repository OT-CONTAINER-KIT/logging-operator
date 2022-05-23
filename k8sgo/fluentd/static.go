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

package k8sfluentd

// configMapContent is static content for fluentd
const configMapContent = `@include /fluentd/etc/conf.d/*.conf
<match fluent.**>
@type null
</match>
<source>
@type tail
path /var/log/containers/*.log
exclude_path ["/var/log/containers/*kube-system*.log", "/var/log/containers/*monitoring*.log", "/var/log/containers/*logging*.log"]
pos_file /var/log/fluentd-containers.log.pos
time_format %Y-%m-%dT%H:%M:%S.%NZ
tag kubernetes.*
format json
read_from_head false
</source>
<filter kubernetes.**>
@type kubernetes_metadata
verify_ssl false
</filter>
<filter kubernetes.**>
@type parser
key_name log
reserve_time true
reserve_data true
emit_invalid_record_to_error false
format json
<parse>
  @type json
</parse>
</filter>
<match kubernetes.**>
  @type elasticsearch_dynamic
  include_tag_key true
  logstash_format true
  # logstash_prefix kubernetes-${record['kubernetes']['pod_name']}
  logstash_prefix kubernetes-${record['kubernetes']['#{ENV['FLUENT_INDEX_PATTERN']}']}
  host "#{ENV['FLUENT_ELASTICSEARCH_HOST']}"
  port "#{ENV['FLUENT_ELASTICSEARCH_PORT']}"
  scheme "#{ENV['FLUENT_ELASTICSEARCH_SCHEME'] || 'http'}"
  ssl_verify "#{ENV['FLUENT_ELASTICSEARCH_SSL_VERIFY'] || 'true'}"
  user "#{ENV['FLUENT_ELASTICSEARCH_USER']}"
  password "#{ENV['FLUENT_ELASTICSEARCH_PASSWORD']}"
  reload_connections false
  reconnect_on_error true
  reload_on_failure true
  <buffer>
	  flush_thread_count 16
	  flush_interval 5s
	  chunk_limit_size 2M
	  queue_limit_length 32
	  retry_max_interval 30
	  retry_forever true
  </buffer>
</match>
`
