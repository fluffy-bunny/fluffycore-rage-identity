{{define "oidc/htmx/_partials/error"}}
<h2 style="color:#e53e3e;">{{ call .LocalizeMessage "error" }}</h2>

{{ if .errorMessage }}
<div class="error-message">
    <span>{{ .errorMessage }}</span>
</div>
{{ end }}

{{ if .errorCode }}
<p class="sub-text">Code: {{ .errorCode }}</p>
{{ end }}

<div class="button-group">
    <button type="button" class="btn-primary"
            hx-get="{{ .paths.HTMXHome }}"
            hx-target="#main-content"
            hx-swap="innerHTML">
        {{ call .LocalizeMessage "start_over" }}
    </button>
</div>
{{end}}
