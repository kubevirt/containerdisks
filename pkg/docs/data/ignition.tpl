{
  "ignition": {
    "version": "3.3.0"
  },
  "passwd": {
    "users": [
      {
        "name": "{{ or .Username "admin"}}",
        "sshAuthorizedKeys": [
          {{if len .AuthorizedKeys -}}
          {{Join (Quote .AuthorizedKeys) ","}}
          {{- else -}}
          "ssh-rsa AAAA..."
          {{- end}}
        ]
      }
    ]
  }
}