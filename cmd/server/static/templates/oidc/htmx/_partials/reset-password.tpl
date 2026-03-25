{{define "oidc/htmx/_partials/reset-password"}}
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
                <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "password_reset" }}</h2>
                <form hx-post="{{ .paths.HTMXResetPassword }}"
                      hx-target="#main-content"
                      hx-swap="innerHTML"
                      hx-indicator="#reset-indicator">
                    <input type="hidden" name="csrf" value="{{ .csrf }}">
                    <input type="hidden" name="email" value="{{ .email }}">

                    <div class="mb-3">
                        <label for="password" class="form-label">{{ call .LocalizeMessage "new_password" }}</label>
                        <input type="password" class="form-control" id="password" name="password" required autofocus>
                    </div>
                    <div class="mb-3">
                        <label for="confirmPassword" class="form-label">{{ call .LocalizeMessage "confirm_password" }}</label>
                        <input type="password" class="form-control" id="confirmPassword" name="confirmPassword" required>
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
                            <span id="reset-indicator" class="htmx-indicator spinner-border spinner-border-sm ms-2" role="status"></span>
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
{{end}}
