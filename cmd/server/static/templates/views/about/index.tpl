{{define "views/about/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
  
    <div class="text-center mt-5" class="alert alert-success" role="alert">
        <h1>Perfect Corp.</h1>
        <p class="lead">Everything good, nothing bad</p>
        <div class="mt-5 alert alert-success" class="alert alert-success" role="alert">
            <table class="table table-striped">
                <thead>
                <tr>
                <th class="text-start" scope="col">#</th>
                <th class="text-start" scope="col">Verbs</th>
                <th class="text-start" scope="col">Path</th>
                </tr>
            </thead>
            <tbody>
            {{range $idx,$def := .defs}}
                
                <tr>
                <th class="text-start" scope="row">{{$idx}}</th>
                <td class="text-start">{{$def.Verbs}}</td>
                <td class="text-start">{{$def.Path}}</td>
                </tr>
            {{end}}
            </tbody>
            </table>
        </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}