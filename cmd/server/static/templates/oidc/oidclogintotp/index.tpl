{{define "oidc/oidclogintotp/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $directive   := .directive }}
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
                        <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "totp_authenticator_app_login" }}</h2>
                        <form action="{{ $paths.OIDCLoginTOTP }}" method="post">
                             {{ if .verified }}
                             {{ else }}
                            <div class="mb-3">
                                <!-- Display QR code here (you'll need a library like QRCode.js) -->
                                <img src="data:image/png;base64,{{ .pngQRCode }}" alt="QR Code" style="max-width: 100%; max-height: 100%;" />
                            </div>
                         {{ end }}
                             <div class="mb-3">
                                <label for="username" class="form-label">Email address</label>
                                <input type="email" class="form-control" id="username" name="username" value="{{ .email }}" required readonly>
                            </div>
                            <div class="mb-3">
                                <label for="code" class="form-label">Code</label>
                                <input type="text" class="form-control" id="code" name="code" required>
                            </div>
              
                            <div class="d-flex justify-content-between">
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