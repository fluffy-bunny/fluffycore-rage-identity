{{define "oidc/htmx/_partials/verify-code"}}
{{ if .errors }}
<div class="error-message">
    {{range .errors}}
    <span>{{.}}</span>
    {{end}}
</div>
{{ end }}

<h2>{{ call .LocalizeMessage "verifycode" }}</h2>
<p>A verification code has been emailed to {{ .email }}.</p>

<form hx-post="{{ .paths.HTMXVerifyCode }}"
      hx-target="#main-content"
      hx-swap="innerHTML"
      hx-indicator="#verify-indicator">
    <input type="hidden" name="csrf" value="{{ .csrf }}">
    <input type="hidden" name="email" value="{{ .email }}">
    <input type="hidden" name="directive" value="{{ .directive }}">

    <div class="form-group">
        <label for="code">{{ call .LocalizeMessage "code" }}</label>
        <input type="text" class="verification-input" id="code" name="code" value="{{ .code }}" required autofocus>
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
            <span id="verify-indicator" class="htmx-indicator" role="status"> ...</span>
        </button>
    </div>
</form>
{{end}}
