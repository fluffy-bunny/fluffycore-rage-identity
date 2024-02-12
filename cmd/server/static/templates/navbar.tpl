{{define "navbar"}}
<!-- Responsive navbar-->
<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container">
        <a class="navbar-brand" href="{{ .paths.Home }}">Echo Starter</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation"><span class="navbar-toggler-icon"></span></button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
                <li class="nav-item"><a class="nav-link active" aria-current="page" href="{{ .paths.About }}">About</a></li>
                <li class="nav-item"><a class="nav-link active" aria-current="page" href="{{ .paths.Login }}">Login</a></li>
                <li class="nav-item dropdown">
                    <a class="nav-link dropdown-toggle" id="navbarDropdown" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">{{ .username }}</a>
                    <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="navbarDropdown">
                    {{ if (call .isAuthenticated ) }}
                        <li><a class="dropdown-item" href="{{ .paths.Logout }}">Logout</a></li>
                    {{ else }}
                        <li><a class="dropdown-item" href="{{ .paths.Login }}">Login</a></li>
                    {{end}}
                        <li><a class="dropdown-item" href="{{ .paths.Profile }}">Profile</a></li>
                    </ul>
                </li>

            </ul>
        </div>
    </div>
</nav>
{{end}}