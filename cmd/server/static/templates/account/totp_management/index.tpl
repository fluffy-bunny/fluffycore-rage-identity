{{define "account/totp_management/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
{{ $returnUrl   := .returnUrl }}

<body>
    <!-- Page content -->
    <div class="container">
        <div class="text-center mt-5">
            <h1>{{ call .LocalizeMessage "totp_management" }}</h1>
         </div>
         {{ if .verified }}
            <form action="{{ .formAction }}" method="post">
   
                <input type="hidden" name="returnUrl" value="{{ $returnUrl }}">
                <input type="hidden" name="action" value="unenroll">
                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "totp_unenroll" }}</button>
            </form>
         {{ else }}
            <form action="{{ .formAction }}" method="post">
                <div class="mb-3">
                    <!-- Display QR code here (you'll need a library like QRCode.js) -->
                    <img src="data:image/png;base64,{{ .pngQRCode }}" alt="QR Code" style="max-width: 100%; max-height: 100%;" />
                </div>
                <input type="hidden" name="returnUrl" value="{{ $returnUrl }}">
                <input type="hidden" name="action" value="enroll">
                <div class="mb-3">
                    <label for="code" class="form-label">{{ call .LocalizeMessage "code" }}</label>
                    <input type="code" class="form-control" id="code" name="code" placeholder="{{ call .LocalizeMessage "totp_enter_placeholder" }}" required>
                </div>
                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "totp_enroll" }}</button>
            </form>
         {{ end }}
    </div>
</body>


{{template "footer" .}}
{{template "html_end" .}}
{{end}}
 