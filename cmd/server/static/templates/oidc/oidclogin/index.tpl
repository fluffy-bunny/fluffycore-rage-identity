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
                <div class="error-container">
                    <ul class="error-list">
                        {{range $idx, $error := .errors}}
                            <li class="error-list-item">
                                {{$error}}
                            </li>
                        {{end}}
                    </ul>
                </div>
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
<script>
const url = '/api'; // relative URL
const data = { request_type: 'InitialPageRequest', version: '1' };  

fetch(url, {
  method: 'POST', 
  mode: 'cors', 
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(data) 
})
.then(response => response.json())
.then(data => console.log(data))
.catch((error) => {
  console.error('Error:', error);
});

</script>
{{template "html_end" .}}
{{end}}