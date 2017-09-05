# Certdump with pkcs8 support

Based off [certdump](https://gist.github.com/tam7t/1b45125ae4de13b3fc6fd0455954c08e) consul-template plugin for writing vault-generated certificates to separate files with support for [PKCS8](https://en.wikipedia.org/wiki/PKCS_8) private key format.

# Usage:
```
certdump <filepath> <user> <group> <mode> <type> <data>
```

`type` is optional, supported: none(default), pkcs8

Example of use in consul-template:
```
{{ with secret "pki/issue/logstash" "ttl=720h" "common_name=logstash.service.consul" }}
{{ .Data.serial_number }}
{{ .Data.certificate | plugin "certdump" "/srv/ssl/logstash.pem" "logstash" "root" "0444" }}
{{ .Data.private_key | plugin "certdump" "/srv/ssl/logstash-key.pem" "logstash" "root" "0400" }}
{{ .Data.private_key | plugin "certdump" "/srv/ssl/logstash-key-pkcs8.pem" "logstash" "root" "0400" "pkcs8" }}
{{ .Data.issuing_ca | plugin "certdump" "/srv/ssl/logstash-ca.pem" "root" "root" "0444" }}
{{ end }}
```

# `certdump` build:
```
go get -v -t -d ./...
go build certdump.go
```
