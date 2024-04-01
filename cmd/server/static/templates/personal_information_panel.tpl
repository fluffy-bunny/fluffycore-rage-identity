
{{define "personal_information_panel"}}
{{ $paths       := .paths }}
{{ $returnUrl   := .returnUrl }}
{{ $readonly    := "" }}
{{ if .displayOnly }}
    {{$readonly = "readonly"}}
{{end}}
{{ $csrf        := .csrf }}

<div class="card shadow">
    <div class="card-body p-4">
        <h2 class="card-title text-center mb-4">{{ call .LocalizeMessage "personal_information" }}</h2>
        <form action="{{ .formAction }}" method="post">
            <input type="hidden" name="csrf" value="{{ $csrf }}">
    
            <input type="hidden" name="action" value="{{ .action }}">
            <input type="hidden" name="returnUrl" value="{{ $returnUrl }}">
            <div class="mb-3">
                <label for="email" class="form-label">{{ call .LocalizeMessage "email" }}</label>
                <input type="email" class="form-control" id="email" name="email" value={{ .email}} required readonly>
            </div>
            <div class="mb-3">
                <label for="given_name" class="form-label">{{ call .LocalizeMessage "given_name" }}</label>
                <input type="text" class="form-control" id="given_name" name="given_name" value="{{ .given_name}}" {{ $readonly }}>
            </div>
            <div class="mb-3">
                <label for="family_name" class="form-label">{{ call .LocalizeMessage "family_name" }}</label>
                <input type="text" class="form-control" id="family_name" name="family_name" value="{{ .family_name}}" {{ $readonly }} >
            </div>
            <div class="mb-3">
                <label for="phone_number" class="form-label">{{ call .LocalizeMessage "phone_number" }}</label>
                <input type="text" class="form-control" id="phone_number" name="phone_number" value="{{ .phone_number}}" {{ $readonly }} >
            </div>
            {{ if .displayOnly }}
                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "edit" }}</button>
            {{ else }}
                <button type="submit" class="btn btn-primary btn-block">{{ call .LocalizeMessage "next" }}</button>
            {{ end }}
        </form>
    </div>
</div>

{{end}}