{{define "account/profile/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}

{{ $paths       := .paths }}
{{ $csrf        := .csrf }}

<body>
    <!-- Page content -->
    <div class="container">
        <div class="text-center mt-5">
            <h1>{{ call .LocalizeMessage "user_profile" }}</h1>
            <!-- Collapsible panels -->
            <div class="accordion" id="linkAccordion">
                <div class="accordion-item">
                    <h2 class="accordion-header" id="personalInformation">
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
                    <h2 class="accordion-header" id="passwordManagement">
                        <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOpenID" aria-expanded="false" aria-controls="collapseOpenID">
                             {{ call .LocalizeMessage "security_settings" }}
                        </button>
                    </h2>
                    <div id="collapseOpenID" class="accordion-collapse collapse" aria-labelledby="passwordManagement" data-bs-parent="#linkAccordion">
                        <div class="accordion-body">
                            <form action="{{ $paths.Profile }}" method="post">
                                <input type="hidden" name="action" value="password-reset">
                                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "password_reset" }}</button>
                            </form>
                            <form action="{{ $paths.Profile }}" method="post">
                                <input type="hidden" name="action" value="totp-management">
                                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "totp_management" }}</button>
                            </form>
                        </div>
                    </div>
                </div>

                <div class="accordion-item">
                    <h2 class="accordion-header" id="passKeysManagement">
                        <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapsePassKeys" aria-expanded="false" aria-controls="collapsePassKeys">
                             {{ call .LocalizeMessage "pass_keys" }}
                        </button>
                    </h2>
                    <div id="collapsePassKeys" class="accordion-collapse collapse" aria-labelledby="passKeysManagement" data-bs-parent="#linkAccordion">
                        <div class="accordion-body">
                            <form action="{{ $paths.Profile }}" method="post">
                                <input type="hidden" name="action" value="passkeys">
                                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "passkey_management" }}</button>
                            </form>
                            {{ if .user.RageUser.WebAuthN }}
                            <div class="error-container">
                                <ul class="error-list">
                                    {{range $idx, $credential := .user.RageUser.WebAuthN.Credentials}}
                                        <li class="error-list-item">
                                            {{ $credential.Authenticator.FriendlyName }}
                                        </li>
                                    {{end}}
                                </ul>
                            </div>
                            {{ end }}
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
