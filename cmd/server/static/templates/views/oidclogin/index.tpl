{{define "views/oidclogin/index"}}
{{template "html_begin" .}}
{{template "header" .}}

{{ $state       := .state }}
{{ $directive   := .directive }}

<body>
<!-- Page content-->
<div class="container">
   
    <div class="text-center mt-5" class="alert alert-success" role="alert">
        <h1>{{ .login }}</h1>
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
        <form action="/oidc-login" method="post">
            <input type="hidden" name="state" value="{{ $state  }}">
            <div class="mb-3">
                <label for="username" class="form-label">Username</label>
                <input type="text" class="form-control" id="username" name="username" required>
            </div>
            <div class="mb-3">
                <label for="password" class="form-label">Password</label>
                <input type="password" class="form-control" id="password" name="password" required>
            </div>
            <button type="submit" class="btn btn-primary">{{ call .LocalizeMessage "login" }}</button>
        </form>
        <p><a class="nav-link active" aria-current="page" href="{{ .paths.Signup }}?state={{ $state }}&wizard_mode=true">{{ call .LocalizeMessage "signup" }}</a></p>
        <p><a class="nav-link active" aria-current="page" href="{{ .paths.ForgotPassword }}?state={{ $state }}">{{ call .LocalizeMessage "forgot_password" }}</a></p>

        <div class="text-center mt-5" class="alert alert-success" role="alert">
        {{range $idx,$idp := .idps}}
            <form action="/external-idp" method="post">
                <input type="hidden" name="state"       value="{{ $state }}">
                <input type="hidden" name="directive"   value="{{ $directive }}">
                <input type="hidden" name="idp_hint"    value="{{$idp.Slug}}">
                <button type="submit" class="btn btn-primary">{{$idp.Name}}</button>
            </form>
        {{end}}
    </div>
    </div>
</div>
</body>
    
{{template "footer" .}}
{{template "html_end" .}}
{{end}}