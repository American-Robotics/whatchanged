## {{ .Version }} ({{ .Date }})

{{- define "body" -}}
{{range . -}}
- {{if .Field.Header.Scope }}**{{ unescape .Field.Header.Scope }}**: {{ end }}{{ unescape .Field.Header.Subject }}({{ hashURL .Hash}}) {{if .Field.Footer }} {{if .Field.Footer.Closes }}, Closes: {{ stringsJoin .Field.Footer.Closes "," }} {{- end }}  {{- end}}
{{ end }}
{{- end -}}

{{if .Feat}}
### 🔥  New feature:
{{ template "body" .Feat }}
{{ end }}

{{if .Fix}}
### 🐛  Bugs fixed:
{{ template "body" .Fix }}
{{ end }}

{{if .Revert}}
### 🔙 Revert:
{{range .Revert -}}
- {{if .RevertCommitHash }}revert {{ hashURL .RevertCommitHash }}, {{ end }}{{ unescape .Field.Header.Subject }}({{ hashURL .Hash}})
{{ end }}
{{ end }}

{{if .BreakingChanges}}
### ❤️ BREAKING CHANGES:
{{ range .BreakingChanges -}}

- {{if .Field.Footer.BreakingChange.Title}}{{ unescape .Field.Footer.BreakingChange.Title }}{{ else }}{{ unescape .Field.Title }}{{ end }}

{{ unescape .Field.Footer.BreakingChange.Content }}

{{ end -}}
{{ end }}
