{{define "oidc/htmx/_partials/forgot-password"}}
<div class="row justify-content-center">
    <div class="col-md-6">
        {{ if .errors }}
        <div class="alert alert-danger" role="alert">
            <ul class="mb-0">
                {{range .errors}}
                <li>{{.}}</li>
                {{end}}
            </ul>
        </div>
        {{ end }}
        <div class="card shadow">
            <div class="card-body p-4">
                <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "forgot_password" }}</h2>
                <p class="text-muted">{{ call .LocalizeMessage "enter_email_for_reset" }}</p>
                <form hx-post="{{ .paths.HTMXForgotPassword }}"
                      hx-target="#main-content"
                      hx-swap="innerHTML"
                      hx-indicator="#forgot-indicator">
                    <input type="hidden" name="csrf" value="{{ .csrf }}">

                    <div class="mb-3">
                        <label for="email" class="form-label">{{ call .LocalizeMessage "email" }}</label>
                        <input type="email" class="form-control" id="email" name="email" value="{{ .email }}" required autofocus>
                    </div>

                    <div class="d-flex justify-content-between">
                        <button type="button" class="btn btn-outline-secondary"
                                hx-get="{{ .paths.HTMXHome }}"
                                hx-target="#main-content"
                                hx-swap="innerHTML">
                            {{ call .LocalizeMessage "cancel" }}
                        </button>
                        <button type="submit" class="btn btn-primary">
                            {{ call .LocalizeMessage "next" }}
                            <span id="forgot-indicator" class="htmx-indicator spinner-border spinner-border-sm ms-2" role="status"></span>
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
{{end}}
