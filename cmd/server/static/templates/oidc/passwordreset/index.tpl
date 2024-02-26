{{define "oidc/passwordreset/index"}}
{{template "html_begin" .}}
{{template "header" .}}

 {{ $paths       := .paths }}

<body class="bg-light d-flex align-items-center min-vh-100">
    <div class="container">
    
        <div class="row justify-content-center">
            <div class="col-md-6">
                {{ if len .errors }}
                
                <table class="table table-bordered">
                    <thead>
                        <tr>
                            <th scope="col">ID</th>
                            <th scope="col">Name</th>
                            <th scope="col">Message</th>
                        </tr>
                    </thead>
                    <tbody>
                    {{range $idx,$def := .errors}}
                        <tr>
                            <th class="text-start" scope="row">{{$idx}}</th>
                            <td class="text-start">{{$def.Key}}</td>
                            <td class="text-start">{{$def.Value}}</td>
                        </tr>
                        
                        {{end}}
                    </tbody>
                </table>
          
                {{ end }}
                {{template "password_reset_panel" .}}
            </div>
        </div>
    </div>
 
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}