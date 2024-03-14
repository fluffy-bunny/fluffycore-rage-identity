{{define "oidc/oidcloginpasskey/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $directive   := .directive }}
{{ $paths       := .paths }}

<script>
window.onload = function() {
    LoginUser();
};
</script>

{{template "footer" .}}
<script src="/static/js/webauthn.js"></script>

 
{{template "html_end" .}}
{{end}}