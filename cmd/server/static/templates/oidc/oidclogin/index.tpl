{{define "oidc/oidclogin/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $directive   := .directive }}
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
                    {{range $idx,$error := .errors}}
                        <tr>
                            <th class="text-start" scope="row">{{$idx}}</th>
                            <td class="text-start">{{$error.Key}}</td>
                            <td class="text-start">{{$error.Value}}</td>
                        </tr>
                        
                        {{end}}
                    </tbody>
                </table>
          
                {{ end }}
                <div class="card shadow">
                    <div class="card-body p-4">
                        <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "login" }}</h2>
                        <form action="{{ $paths.OIDCLogin }}" method="post">
                             <div class="mb-3">
                                <label for="username" class="form-label">Email address</label>
                                <input type="email" class="form-control" id="username" name="username" placeholder="Enter your email" value="{{ .email }}" required>
                            </div>
              
                            <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "next" }}</button>
                        </form>
                        <p class="mt-0 text-center">
                            <form action="{{ $paths.ForgotPassword }}" method="post">
                                 <input type="hidden" name="type"       value="GET">          
                                <button type="submit" class="btn btn-link text-muted">{{ call .LocalizeMessage "forgot_password" }}</button>
                            </form>
                            <form action="{{ $paths.Signup }}" method="post">
                                 <input type="hidden" name="type"       value="GET">          
                                <button type="submit" class="btn btn-link text-muted">{{ call .LocalizeMessage "signup" }}</button>
                            </form>    
                         </p>
                         <hr>
                        <p class="text-center">{{ call .LocalizeMessage "or_signin_with" }}</p>
                        <div class="d-flex justify-content-center">
                            {{range $idx,$idp := .idps}}
                                <form action="{{ $paths.ExternalIDP }}" method="post">
                                     <input type="hidden" name="directive"   value="{{ $directive }}">
                                    <input type="hidden" name="idp_hint"    value="{{$idp.Slug}}">
                                    <button type="submit" class="btn btn-outline-primary me-2 ">{{$idp.Name}}</button>
                                </form>
                            {{end}}
                          
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
 
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}