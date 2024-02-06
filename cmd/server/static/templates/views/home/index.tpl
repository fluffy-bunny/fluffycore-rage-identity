{{define "views/home/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>A Bootstrap 5 Starter Template</h1>
        <p class="lead">A complete project boilerplate built with Bootstrap</p>
        <p>Bootstrap v5.1.3</p>
        <p><a class="nav-link active" aria-current="page" href="/.well-known/openid-configuration">openid-configuration</a></p>
        <p><a class="nav-link active" aria-current="page" href="/.well-known/jwks">JWKS</a></p>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}