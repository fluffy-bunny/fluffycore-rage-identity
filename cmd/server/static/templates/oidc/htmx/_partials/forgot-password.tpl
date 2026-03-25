{{define "oidc/htmx/_partials/forgot-password"}}
{{ if .errors }}
<div class="error-message">
    {{range .errors}}
    <span>{{.}}</span>
    {{end}}
</div>
{{ end }}

<h2>{{ call .LocalizeMessage "forgot_password" }}</h2>
<p>{{ call .LocalizeMessage "enter_email_for_reset" }}</p>

<form hx-post="{{ .paths.HTMXForgotPassword }}"
      hx-target="#main-content"
      hx-swap="innerHTML"
      hx-indicator="#forgot-indicator">
    <input type="hidden" name="csrf" value="{{ .csrf }}">

    <div class="form-group">
        <label for="email">{{ call .LocalizeMessage "email" }}</label>
        <input type="email" id="email" name="email" value="{{ .email }}" required autofocus>
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
            <span id="forgot-indicator" class="htmx-indicator" role="status"> ...</span>
        </button>
    </div>
</form>
{{end}}
