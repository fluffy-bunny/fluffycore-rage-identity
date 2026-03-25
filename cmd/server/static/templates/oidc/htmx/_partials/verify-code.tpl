{{define "oidc/htmx/_partials/verify-code"}}
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
                <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "verifycode" }}</h2>
                <p class="text-muted">A verification code has been emailed to {{ .email }}.</p>
                <form hx-post="{{ .paths.HTMXVerifyCode }}"
                      hx-target="#main-content"
                      hx-swap="innerHTML"
                      hx-indicator="#verify-indicator">
                    <input type="hidden" name="csrf" value="{{ .csrf }}">
                    <input type="hidden" name="email" value="{{ .email }}">
                    <input type="hidden" name="directive" value="{{ .directive }}">

                    <div class="mb-3">
                        <label for="code" class="form-label">{{ call .LocalizeMessage "code" }}</label>
                        <input type="text" class="form-control" id="code" name="code" value="{{ .code }}" required autofocus>
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
                            <span id="verify-indicator" class="htmx-indicator spinner-border spinner-border-sm ms-2" role="status"></span>
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
{{end}}
