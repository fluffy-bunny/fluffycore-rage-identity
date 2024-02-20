{{define "oidc/forgotpassword/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $state       := .state }}

<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5" class="alert alert-success" role="alert">
        <h1> {{ call .LocalizeMessage "forgotpassword" }}</h1>
        <div class="mt-5 alert alert-success" class="alert alert-success" role="alert">
            <table class="table table-striped">
                <thead>
                <tr>
                <th class="text-start" scope="col">#</th>
                <th class="text-start" scope="col">Key</th>
                <th class="text-start" scope="col">Value</th>
                </tr>
            </thead>
            <tbody>

            {{range $idx,$def := .defs}}
                <tr>
                <th class="text-start" scope="row">{{$idx}}</th>
                <td class="text-start">{{$def.Key}}</td>
                <td class="text-start">{{$def.Value}}</td>
                </tr>
            {{end}}
            </tbody>
            </table>
        </div>
        <form action="{{ .paths.ForgotPassword }}" method="post">
            <input type="hidden" name="state" value="{{ $state  }}">
            <div class="mb-3">
                <label for="email" class="form-label">Email</label>
                <input type="text" class="form-control" id="email" name="email" required>
            </div>       
            <button type="submit" class="btn btn-primary">Submit</button>
        </form>
    </div>
</div>
<script>
    // Set the default email value
    document.getElementById('email').value = '{{ .email }}';
</script>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}