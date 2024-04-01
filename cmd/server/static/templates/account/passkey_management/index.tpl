{{define "account/passkey_management/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
{{ $csrf        := .csrf }}

<body>
    <!-- Page content -->
    <div class="container">
        <div class="text-center mt-5">
            <h1>{{ call .LocalizeMessage "passkey_management" }}</h1>
            <button class="btn btn-outline-primary" onclick="registerUser({{ .returnUrl }})">{{ call .LocalizeMessage "register" }}</button>
        </div>
    </div>
</body>


{{template "footer" .}}
<script src="/static/js/webauthn.js"></script>

{{template "html_end" .}}
{{end}}
