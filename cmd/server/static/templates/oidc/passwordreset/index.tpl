{{define "oidc/passwordreset/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $paths       := .paths }}
{{ $csrf        := .csrf }}

<body class="bg-light d-flex align-items-center min-vh-100">
    <div class="container">
    
        <div class="row justify-content-center">
            <div class="col-md-6">
                {{ if len .errors }}
                <div class="error-container">
                    <ul class="error-list">
                        {{range $idx, $error := .errors}}
                            <li class="error-list-item">
                                {{$error}}
                            </li>
                        {{end}}
                    </ul>
                </div>
                {{ end }}
                {{template "password_reset_panel" .}}
            </div>
        </div>
    </div>
 
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}