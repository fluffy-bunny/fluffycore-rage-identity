{{define "oidc/htmx/_partials/signup"}}
{{ if .errors }}
<div class="error-message">
    {{range .errors}}
    <span>{{.}}</span>
    {{end}}
</div>
{{ end }}

<h2>{{ call .LocalizeMessage "signup" }}</h2>
<p>Create a new account</p>

<form hx-post="{{ .paths.HTMXSignup }}"
      hx-target="#main-content"
      hx-swap="innerHTML"
      hx-indicator="#signup-indicator">
    <input type="hidden" name="csrf" value="{{ .csrf }}">

    <div class="form-group">
        <label for="email">{{ call .LocalizeMessage "email" }}</label>
        <input type="email" id="email" name="email" value="{{ .email }}" required autofocus>
    </div>
    <div class="form-group">
        <label for="password">{{ call .LocalizeMessage "password" }}</label>
        <input type="password" id="password" name="password" required>
    </div>

    <div class="button-group">
        <button type="button" class="btn-secondary"
                hx-get="{{ .paths.HTMXHome }}"
                hx-target="#main-content"
                hx-swap="innerHTML">
            {{ call .LocalizeMessage "cancel" }}
        </button>
        <button type="submit" class="btn-primary">
            {{ call .LocalizeMessage "next" }}
            <span id="signup-indicator" class="htmx-indicator" role="status"> ...</span>
        </button>
    </div>
</form>

{{ if .socialIdps }}
<div class="social-login-section">
    <div class="divider"><span>{{ call .LocalizeMessage "or_signin_with" }}</span></div>
    <div class="social-logins">
        {{range .socialIdps}}
        <form hx-post="{{ $.paths.HTMXHome }}" hx-target="#main-content" hx-swap="innerHTML">
            <input type="hidden" name="csrf" value="{{ $.csrf }}">
            <input type="hidden" name="idp_hint" value="{{.Slug}}">
            <button type="submit" class="btn-social">{{.Slug}}</button>
        </form>
        {{end}}
    </div>
</div>
{{ end }}
{{end}}
