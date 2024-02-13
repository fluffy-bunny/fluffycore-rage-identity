{{define "views/logout/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
 
<div class="container">
 
</div>
 {{ $url := .url }}
<script>
    function getMyRedirectURL() {
       document.write("It will redirect within 1 seconds.....please wait...");//it will redirect after 3 seconds
       setTimeout(function() {
            window.location = {{$url}};
       }, 1000);
    }
    getMyRedirectURL();
</script>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}