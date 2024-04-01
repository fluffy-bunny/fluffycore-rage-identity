{{define "account/personal_information/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}

{{ $paths       := .paths }}
{{ $csrf        := .csrf }}

<body>
    <!-- Page content -->
    <div class="container">
    <div class="text-center mt-5">
        <h1>{{ call .LocalizeMessage "personal_information" }}</h1>

        <div class="row justify-content-center">
             <div class="col-md-6">
                {{template "personal_information_panel" .}}
            </div>
        </div>
    </div>
    </div>
</body>


{{template "footer" .}}
{{template "html_end" .}}
{{end}}
