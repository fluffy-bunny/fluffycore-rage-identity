{{define "head_links"}}
  {{range $idx,$link := .headLinks}}
    <link href="{{ $link.HREF }}" rel="{{ $link.REL }}" />      
  {{end}}
{{end}}

{{define "header"}}
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    {{template "head_links" .}}

</head>
{{end}}

