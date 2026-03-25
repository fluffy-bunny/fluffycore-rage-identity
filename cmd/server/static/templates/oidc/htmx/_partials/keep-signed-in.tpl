{{define "oidc/htmx/_partials/keep-signed-in"}}
<div class="row justify-content-center">
    <div class="col-md-6">
        <div class="card shadow">
            <div class="card-body p-4">
                <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "keep_signed_in" }}</h2>
                <form hx-post="{{ .paths.HTMXKeepSignedIn }}"
                      hx-target="#main-content"
                      hx-swap="innerHTML"
                      hx-indicator="#ksi-indicator">
                    <input type="hidden" name="csrf" value="{{ .csrf }}">

                    <div class="mb-3 form-check">
                        <input type="checkbox" class="form-check-input" id="keepSignedIn" name="keepSignedIn" value="true">
                        <label class="form-check-label" for="keepSignedIn">{{ call .LocalizeMessage "keep_me_signed_in" }}</label>
                    </div>

                    <div class="d-grid">
                        <button type="submit" class="btn btn-primary">
                            {{ call .LocalizeMessage "continue" }}
                            <span id="ksi-indicator" class="htmx-indicator spinner-border spinner-border-sm ms-2" role="status"></span>
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
{{end}}
