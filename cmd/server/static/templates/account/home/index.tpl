{{define "account/home/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
{{ $csrf        := .csrf }}

<body>
    <!-- Page content -->
    <div class="container">
        <div class="text-center mt-5">
            <h1>{{ call .LocalizeMessage "organization_name" }}</h1>
            <!-- Collapsible panels -->
            <div class="accordion" id="linkAccordion">
                <div class="accordion-item">
                    <h2 class="accordion-header" id="swaggerLink">
                        <button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseSwagger" aria-expanded="true" aria-controls="collapseSwagger">
                            Swagger
                        </button>
                    </h2>
                    <div id="collapseSwagger" class="accordion-collapse collapse show" aria-labelledby="swaggerLink" data-bs-parent="#linkAccordion">
                        <div class="accordion-body">
                            <a class="nav-link active" aria-current="page" href="/swagger/">Go to Swagger</a>
                        </div>
                    </div>
                </div>

                <div class="accordion-item">
                    <h2 class="accordion-header" id="openidConfigLink">
                        <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOpenID" aria-expanded="false" aria-controls="collapseOpenID">
                            OpenID Configuration
                        </button>
                    </h2>
                    <div id="collapseOpenID" class="accordion-collapse collapse" aria-labelledby="openidConfigLink" data-bs-parent="#linkAccordion">
                        <div class="accordion-body">
                            <a class="nav-link active" aria-current="page" href="/.well-known/openid-configuration">Go to OpenID Configuration</a>
                        </div>
                    </div>
                </div>

                <div class="accordion-item">
                    <h2 class="accordion-header" id="jwksLink">
                        <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseJWKS" aria-expanded="false" aria-controls="collapseJWKS">
                            JWKS
                        </button>
                    </h2>
                    <div id="collapseJWKS" class="accordion-collapse collapse" aria-labelledby="jwksLink" data-bs-parent="#linkAccordion">
                        <div class="accordion-body">
                            <a class="nav-link active" aria-current="page" href="/.well-known/jwks">Go to JWKS</a>
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