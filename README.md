# bosh-evaluator

```
bosh-evaluator fetch config.yml
```

```
defaults:
  vault:
    path: "secret/redis-((foundation-name))-props"
variables:
- name: "deployment-name"
  value: redis-((foundation-name))
  vault:
    path: "secret/redis-((foundation-name))-props"
- name: "network-name"
  default: "default"
  vault:
    path: "secret/redis-((foundation-name))-props"
- name: "vm-type"
  default: "medium"
  vault:
    path: "secret/redis-((foundation-name))-props"
- name: "disk-type"
  default: "medium"
  vault:
    path: "secret/redis-((foundation-name))"
- name: "broker-ip"
  type: ips
- name: "cf-admin-username"
- name: "cf-admin-password"
- name: "cf-apps-domain"
- name: "cf-system-domain"
- name: "dedicated-nodes-count"
- name: "dedicated-nodes-ips"
  type: ips
- name: "nats-ips"
  type: ips
  external: true
- name: "nats-port"
- name: "nats-username"
- name: "nats-password"
- name: "network-name"
- name: "syslog-aggregator-host"
- name: "syslog-aggregator-port"
  vault_path: secrets/cf-((foundation-name))-props
- name: "az"
  default: "az1"
```

