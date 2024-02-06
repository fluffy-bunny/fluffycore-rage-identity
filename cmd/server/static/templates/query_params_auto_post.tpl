{{define "query_params_auto_post"}}
{{template "html_begin" .}}
<head>
      <meta charset="utf-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
      <meta name="description" content="" />
      <meta name="author" content="" />
      <title>oidc-process</title>
</head>
<body onload="setTimeout(function() { document.frm1.submit() }, 0)">
     {{ $csrf := .security.CSRF }}
     <form  method="post" name="frm1">
        <input type="hidden" name="csrf" value="{{ $csrf }}">
     </form>
</body>
{{template "html_end" .}}
{{end}}