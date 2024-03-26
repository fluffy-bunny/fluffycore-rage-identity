{{define "account/totp_management/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}

<body>
    <!-- Page content -->
    <div class="container">
        <div class="text-center mt-5">
            <h1>{{ call .LocalizeMessage "totp_management" }}</h1>
         </div>
         <div class="row justify-content-center mt-4">
            <div class="col-md-6">
                <!-- Display QR code here (you'll need a library like QRCode.js) -->
                <img src="data:image/png;base64,{{ .pngQRCode }}" alt="QR Code" style="max-width: 100%; max-height: 100%;" />
            </div>
    </div>
</body>


{{template "footer" .}}
{{template "html_end" .}}
{{end}}
 