{{define "oidc/htmx/_partials/home"}}
{{ if .errors }}
<div class="error-message">
    {{range .errors}}
    <span>{{.}}</span>
    {{end}}
</div>
{{ end }}

<h2>{{ call .LocalizeMessage "signin" }}</h2>
<p>Enter your email to get started</p>

<form hx-post="{{ .paths.HTMXHome }}"
      hx-target="#main-content"
      hx-swap="innerHTML"
      hx-indicator="#home-indicator">
    <input type="hidden" name="csrf" value="{{ .csrf }}">

    <div class="form-group">
        <label for="email">{{ call .LocalizeMessage "email" }}</label>
        <input type="email" id="email" name="email" value="{{ .email }}" required autofocus>
    </div>

    <div class="button-group">
        <button type="submit" class="btn-primary">
            {{ call .LocalizeMessage "next" }}
            <span id="home-indicator" class="htmx-indicator" role="status"> ...</span>
        </button>
    </div>
</form>

{{ if not .disableSignup }}
<div class="create-account">
    <a href="#" hx-get="{{ .paths.HTMXSignup }}" hx-target="#main-content" hx-swap="innerHTML" hx-push-url="true">
        {{ call .LocalizeMessage "signup" }}
    </a>
</div>
{{ end }}

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
