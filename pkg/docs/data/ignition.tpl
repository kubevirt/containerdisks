{
  "ignition": {
    "version": "3.3.0"
  },
  "passwd": {
    "users": [
      {
        "name": "{{ .Username }}",
        "sshAuthorizedKeys": [
          {{if len .AuthorizedKeys -}}
          {{Join (Quote .AuthorizedKeys) ",\n          "}}
          {{- else -}}
          "ssh-rsa AAAA..."
          {{- end}}
        ]
      }
    ]
  }
}