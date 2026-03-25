{{define "oidc/htmx/_partials/reset-password"}}
{{ if .errors }}
<div class="error-message">
    {{range .errors}}
    <span>{{.}}</span>
    {{end}}
</div>
{{ end }}

<h2>{{ call .LocalizeMessage "password_reset" }}</h2>
<p>Choose a new password for your account</p>

<form hx-post="{{ .paths.HTMXResetPassword }}"
      hx-target="#main-content"
      hx-swap="innerHTML"
      hx-indicator="#reset-indicator">
    <input type="hidden" name="csrf" value="{{ .csrf }}">
    <input type="hidden" name="email" value="{{ .email }}">

    <div class="form-group">
        <label for="password">{{ call .LocalizeMessage "new_password" }}</label>
        <input type="password" id="password" name="password" required autofocus>
    </div>
    <div class="form-group">
        <label for="confirmPassword">{{ call .LocalizeMessage "confirm_password" }}</label>
        <input type="password" id="confirmPassword" name="confirmPassword" required>
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
            <span id="reset-indicator" class="htmx-indicator" role="status"> ...</span>
        </button>
    </div>
</form>
{{end}}
