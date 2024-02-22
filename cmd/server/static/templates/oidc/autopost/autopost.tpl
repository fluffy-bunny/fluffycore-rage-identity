{{define "oidc/autopost/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $form_params := .form_params }}

<body>

<form id="autoForm" action="{{ .action }}" method="post">
    {{range $idx,$formParam := .form_params}}
    <input type="hidden" name="{{ $formParam.Name }}" value="{{ $formParam.Value }}">
    {{ end }}
    <input type="submit" value="Submit">
</form>
<script>
    // Automatically submit the form when the page loads
    document.getElementById("autoForm").submit();
</script>
</body>
{{template "html_end" .}}
{{end}}