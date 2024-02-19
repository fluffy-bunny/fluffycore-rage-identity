{{define "base_begin"}}
{{template "html_begin" .}}
{{template "header" .}}
<body style="font-family:Arial,sans-serif;font-size:14px;padding:10px;margin-bottom: 20px;">
<div style="margin-bottom:10px; text-align: center; padding: 10px; background-color: #80adbe;">
    <p style="font-size: 24px; color: white;margin-top: 0;margin-bottom: 0;">
        {{ .organization }}
    </p>
</div>
<div class="body">
{{end}}
 
{{define "base_end"}}
</div>
<div class="footer" style="color:#444;font-size:12px;margin-top:20px;">
    <hr />
    You received this email because you have a {{ .organization }} account. Manage your account at
    <a href="{{.account_url}}">{{ call .LocalizeMessage "my_account" }}</a>. You can also reach out to <a
        href="mailto:rage@test.com">{{ call .LocalizeMessage "admin" }}</a> for no help at all.</div>
</div>
</body>
{{template "footer" .}}
{{template "html_end" .}}
{{end}}

   
