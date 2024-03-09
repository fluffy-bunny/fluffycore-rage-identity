
{{define "password_reset_panel"}}
 {{ $paths       := .paths }}
{{ $returnUrl   := .returnUrl }}

<div class="card shadow">
    <div class="card-body p-4">
        <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "password_reset" }}</h2>
        <form action="{{ $paths.PasswordReset }}" method="post">
             <input type="hidden" name="returnUrl" value="{{ $returnUrl }}">

            <div class="mb-3">
                <label for="password" class="form-label">{{ call .LocalizeMessage "password" }}</label>
                <input type="password" class="form-control" id="password" name="password" required >
            </div>
            <div class="mb-3">
                <label for="confirmPassword" class="form-label">{{ call .LocalizeMessage "confirm_password" }}</label>
                <input type="password" class="form-control" id="confirmPassword" name="confirmPassword" required >
            </div>
            <div class="d-flex justify-content-between">
                <button type="submit" class="btn btn-outline-primary" name="action" value="cancel" formnovalidate>{{ call .LocalizeMessage "cancel" }}</button>
                <div class="btn-group">
                    <button type="submit" class="btn btn-primary" name="action" value="next">{{ call .LocalizeMessage "next" }}</button>
                </div>
            </div>                        
        </form>
    </div>
</div>
{{end}}