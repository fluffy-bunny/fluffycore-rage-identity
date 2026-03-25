{{define "oidc/htmx/shell"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
    <meta name="description" content="OIDC Login" />
    <title>Login</title>
    <link rel="icon" type="image/x-icon" href="/static/assets/favicon.ico" />
    <link href="/static/go-app/oidc-login/static_output/web/styles.css" rel="stylesheet" />
    <script src="https://unpkg.com/htmx.org@2.0.4" integrity="sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+" crossorigin="anonymous"></script>
    <meta name="htmx-config" content='{"responseHandling":[{"code":".*", "swap": true}]}'>
    <style>
        .htmx-indicator {
            display: none;
        }
        .htmx-request .htmx-indicator,
        .htmx-request.htmx-indicator {
            display: inline-block;
        }
    </style>
</head>
<body>
    <div class="wizard-container">
        <div class="app-header">
            <div class="header-content">
                <div class="logo-title-group">
                    <img src="/static/go-app/oidc-login/static_output/web/m_logo.svg" alt="Logo" class="app-logo" />
                    <div class="title-version-group">
                        <div class="app-title">{{ .brandTitle }}</div>
                    </div>
                </div>
            </div>
        </div>
        <div id="main-content" class="step-container"
             hx-get="{{ .paths.HTMXHome }}"
             hx-trigger="load"
             hx-swap="innerHTML">
            <p style="text-align:center;padding:40px 0;">Loading...</p>
        </div>
    </div>
</body>
</html>
{{end}}
