{{define "oidc/oidcloginpasskey/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $directive   := .directive }}
{{ $paths       := .paths }}
{{ $csrf        := .csrf }}

<script>
window.onload = function() {
    LoginUser({{ .returnFailedUrl }});
};
</script>

{{template "footer" .}}
<script src="/static/js/webauthn.js"></script>

 
{{template "html_end" .}}
{{end}}