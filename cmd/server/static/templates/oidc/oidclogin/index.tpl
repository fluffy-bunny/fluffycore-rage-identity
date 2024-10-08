{{define "oidc/oidclogin/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $directive   := .directive }}
{{ $paths       := .paths }}
{{ $csrf        := .csrf }}

<body class="bg-light d-flex align-items-center min-vh-100">
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <div id="root"></div>
</body>

{{template "footer" .}}
    <script defer="defer" src="/static/oidc-flows/static/js/main.js"></script>

{{template "html_end" .}}
{{end}}