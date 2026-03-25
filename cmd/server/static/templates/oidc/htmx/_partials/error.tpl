{{define "oidc/htmx/_partials/error"}}
<div class="row justify-content-center">
    <div class="col-md-6">
        <div class="card shadow border-danger">
            <div class="card-body p-4">
                <h2 class="card-title text-center text-danger mb-4">{{ call .LocalizeMessage "error" }}</h2>
                {{ if .errorMessage }}
                <div class="alert alert-danger" role="alert">
                    {{ .errorMessage }}
                </div>
                {{ end }}
                {{ if .errorCode }}
                <p class="text-muted text-center"><small>Code: {{ .errorCode }}</small></p>
                {{ end }}
                <div class="d-grid">
                    <button type="button" class="btn btn-primary"
                            hx-get="{{ .paths.HTMXHome }}"
                            hx-target="#main-content"
                            hx-swap="innerHTML">
                        {{ call .LocalizeMessage "start_over" }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>
{{end}}
