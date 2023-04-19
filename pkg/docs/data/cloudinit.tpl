#cloud-config
# The default username is: {{ .Username }}
ssh_authorized_keys:
  {{- range .AuthorizedKeys}}
  - {{.}}
  {{- else }}
  - ssh-rsa AAAA...
  {{- end}}