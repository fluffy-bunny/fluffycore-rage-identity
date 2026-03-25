{{define "oidc/htmx/_partials/home"}}
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
                <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "signin" }}</h2>
                <form hx-post="{{ .paths.HTMXHome }}"
                      hx-target="#main-content"
                      hx-swap="innerHTML"
                      hx-indicator="#home-indicator">
                    <input type="hidden" name="csrf" value="{{ .csrf }}">

                    <div class="mb-3">
                        <label for="email" class="form-label">{{ call .LocalizeMessage "email" }}</label>
                        <input type="email" class="form-control" id="email" name="email" value="{{ .email }}" required autofocus>
                    </div>

                    <div class="d-grid">
                        <button type="submit" class="btn btn-primary">
                            {{ call .LocalizeMessage "next" }}
                            <span id="home-indicator" class="htmx-indicator spinner-border spinner-border-sm ms-2" role="status"></span>
                        </button>
                    </div>
                </form>

                {{ if not .disableSignup }}
                <hr>
                <p class="text-center mb-2">
                    <a href="#" hx-get="{{ .paths.HTMXSignup }}" hx-target="#main-content" hx-swap="innerHTML" hx-push-url="true">
                        {{ call .LocalizeMessage "signup" }}
                    </a>
                </p>
                {{ end }}

                {{ if .socialIdps }}
                <hr>
                <p class="text-center">{{ call .LocalizeMessage "or_signin_with" }}</p>
                <div class="d-flex justify-content-center flex-wrap gap-2">
                    {{range .socialIdps}}
                    <form hx-post="{{ $.paths.HTMXHome }}" hx-target="#main-content" hx-swap="innerHTML">
                        <input type="hidden" name="csrf" value="{{ $.csrf }}">
                        <input type="hidden" name="idp_hint" value="{{.Slug}}">
                        <button type="submit" class="btn btn-outline-primary">{{.Slug}}</button>
                    </form>
                    {{end}}
                </div>
                {{ end }}
            </div>
        </div>
    </div>
</div>
{{end}}
