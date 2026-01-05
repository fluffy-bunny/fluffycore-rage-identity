{{define "oidc/error/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}

<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Error</h1>
        <div class="alert alert-danger" role="alert">
            {{ if .message }}
                <div><strong>{{ .message }}</strong></div>
            {{ else }}
                <div><strong>An error occurred</strong></div>
            {{ end }}
            {{ if .error }}
                <div class="mt-2"><small>Error code: {{ .error }}</small></div>
            {{ end }}
        </div>
        <div class="mt-4">
            <a href="/" class="btn btn-primary">Return to Home</a>
        </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}