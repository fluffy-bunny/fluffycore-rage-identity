{{define "oidc/autoredirect/index"}}
{{template "html_begin" .}}
{{template "header" .}}
<body>
<!-- Page content-->
 
<div class="container">
 
</div>
 {{ $url := .url }}
<script>
    function getMyRedirectURL() {
       setTimeout(function() {
            window.location = {{$url}};
       }, 100);
    }
    getMyRedirectURL();
</script>
</body>
{{template "html_end" .}}
{{end}}