{{define "oidc/htmx/_partials/password"}}
{{ if .errors }}
<div class="error-message">
    {{range .errors}}
    <span>{{.}}</span>
    {{end}}
</div>
{{ end }}

<h2>{{ call .LocalizeMessage "password" }}</h2>
<p>{{ .email }}</p>

<form hx-post="{{ .paths.HTMXPassword }}"
      hx-target="#main-content"
      hx-swap="innerHTML"
      hx-indicator="#password-indicator">
    <input type="hidden" name="csrf" value="{{ .csrf }}">
    <input type="hidden" name="email" value="{{ .email }}">

    <div class="form-group">
        <label for="password">{{ call .LocalizeMessage "password" }}</label>
        <input type="password" id="password" name="password" required autofocus>
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
            <span id="password-indicator" class="htmx-indicator" role="status"> ...</span>
        </button>
    </div>
</form>

<div class="forgot-password">
    <a href="#" hx-get="{{ .paths.HTMXForgotPassword }}" hx-target="#main-content" hx-swap="innerHTML" hx-push-url="true">
        {{ call .LocalizeMessage "forgot_password" }}
    </a>
</div>
{{end}}
