{{define "account/profile/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}

{{ $paths       := .paths }}
<body>
    <!-- Page content -->
    <div class="container">
        <div class="text-center mt-5">
            <h1>{{ call .LocalizeMessage "user_profile" }}</h1>
            <!-- Collapsible panels -->
            <div class="accordion" id="linkAccordion">
                <div class="accordion-item">
                    <h2 class="accordion-header" id="swaggerLink">
                        <button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseSwagger" aria-expanded="true" aria-controls="collapseSwagger">
                            {{ call .LocalizeMessage "personal_information" }}
                        </button>
                    </h2>
                    <div id="collapseSwagger" class="accordion-collapse collapse show" aria-labelledby="swaggerLink" data-bs-parent="#linkAccordion">
                        <div class="accordion-body">
                            <div class="row justify-content-center">
                                <div class="col-md-6">
                                {{template "personal_information_panel" .}}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="accordion-item">
                    <h2 class="accordion-header" id="openidConfigLink">
                        <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOpenID" aria-expanded="false" aria-controls="collapseOpenID">
                             {{ call .LocalizeMessage "security_settings" }}
                        </button>
                    </h2>
                    <div id="collapseOpenID" class="accordion-collapse collapse" aria-labelledby="openidConfigLink" data-bs-parent="#linkAccordion">
                        <div class="accordion-body">
                            <form action="{{ $paths.Profile }}" method="post">
                                <input type="hidden" name="action" value="password-reset">
                                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "password_reset" }}</button>
                            </form>
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
