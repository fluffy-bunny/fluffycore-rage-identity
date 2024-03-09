{{define "oidc/verifycode/index"}}
{{template "html_begin" .}}
{{template "header" .}}

 {{ $paths       := .paths }}

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
                <div class="card shadow">
                    <div class="card-body p-4">
                        <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "verifycode" }}</h2>
                        <p>A verification code has be emailed to {{.email}} If an account exists. </p>

                        <form action="{{ $paths.VerifyCode }}" method="post">
                            <input type="hidden" name="email"       value="{{ .email }}">
                            <input type="hidden" name="directive"   value="{{ .directive }}">

                            <div class="mb-3">
                                <label for="code" class="form-label">Code</label>
                                <input type="text" class="form-control" id="code" name="code" value="{{ .code }}" required>
                            </div>
              
                            <div class="d-flex justify-content-between">
                                <button type="submit" class="btn btn-outline-primary" name="action" value="cancel" formnovalidate>{{ call .LocalizeMessage "cancel" }}</button>
                                <div class="btn-group">
                                    <button type="submit" class="btn btn-primary" name="action" value="next">{{ call .LocalizeMessage "next" }}</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
 
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}