#cloud-config
users:
  - name: {{ or .Username "myuser" }}
    sudo: ALL=(ALL) NOPASSWD:ALL
    ssh_authorized_keys:
      {{- range .AuthorizedKeys}}
      - {{.}}
      {{- else }}
      - ssh-rsa AAAA...
      {{- end}}
