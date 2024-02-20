{{define "emails/verifycode/index"}}
{{template "base_begin" .}}
 
<!-- Page content-->
<div class="container">
  {{ .verification_code_message}}
</div>
{{template "base_end" .}}
{{end}}