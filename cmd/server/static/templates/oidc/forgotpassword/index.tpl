{{define "oidc/forgotpassword/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $state       := .state }}
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
                <div class="card shadow">
                    <div class="card-body p-4">
                        <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "forgot_password" }}</h2>
                        <form action="{{ $paths.ForgotPassword }}" method="post">
                            <input type="hidden" name="state" value="{{ $state }}">
                            <div class="mb-3">
                                <label for="email" class="form-label">Email address</label>
                                <input type="email" class="form-control" id="email" name="email" placeholder="Enter your email" value="{{ .email }}" required>
                            </div>
              
                            <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "next" }}</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
 
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}