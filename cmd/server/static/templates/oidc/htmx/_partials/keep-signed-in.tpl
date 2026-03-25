{{define "oidc/htmx/_partials/keep-signed-in"}}
<h2>{{ call .LocalizeMessage "keep_signed_in" }}</h2>
<p>Choose whether to stay signed in on this device</p>

<form hx-post="{{ .paths.HTMXKeepSignedIn }}"
      hx-target="#main-content"
      hx-swap="innerHTML"
      hx-indicator="#ksi-indicator">
    <input type="hidden" name="csrf" value="{{ .csrf }}">

    <div class="form-group" style="display:flex;align-items:center;gap:10px;">
        <input type="checkbox" id="keepSignedIn" name="keepSignedIn" value="true" style="width:20px;height:20px;">
        <label for="keepSignedIn" style="margin-bottom:0;">{{ call .LocalizeMessage "keep_me_signed_in" }}</label>
    </div>

    <div class="button-group">
        <button type="submit" class="btn-primary">
            {{ call .LocalizeMessage "continue" }}
            <span id="ksi-indicator" class="htmx-indicator" role="status"> ...</span>
        </button>
    </div>
</form>
{{end}}
