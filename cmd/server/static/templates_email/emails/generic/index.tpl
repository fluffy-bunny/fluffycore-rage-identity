{{define "emails/generic/index"}}
{{template "base_begin" .}}
 
<!-- Page content-->
<div class="container">
  {{ .body }}
</div>
{{template "base_end" .}}
{{end}}