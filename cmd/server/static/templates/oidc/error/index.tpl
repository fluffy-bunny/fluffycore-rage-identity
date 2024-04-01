{{define "oidc/error/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
{{ $csrf        := .csrf }}

<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Error</h1>
        <div class="alert alert-danger" role="alert">
        <div>message:{{ .params.Message }}</div>
    
        </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}