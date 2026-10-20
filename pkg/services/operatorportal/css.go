package operatorportal

// operatorPortalCSS is the entire stylesheet for the authenticated portal
// (everything shellPage wraps, including every resource.Mount screen).
// Plain CSS, no build step, no JS framework — matches this repo's general
// preference for a Go+htmx stack over a JS-heavy admin UI.
const operatorPortalCSS = `
* { box-sizing: border-box; }
body {
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
	margin: 0;
	color: #1a1a1a;
	background: #f5f6f8;
}
.op-topbar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 12px 20px;
	background: #1f2937;
	color: #f9fafb;
}
.op-brand { font-weight: 600; }
.op-topbar-right { display: flex; align-items: center; gap: 12px; }
.op-topbar-right form { margin: 0; }

.op-layout { display: flex; min-height: calc(100vh - 49px); }
.op-sidebar {
	width: 200px;
	flex-shrink: 0;
	background: #ffffff;
	border-right: 1px solid #e2e4e9;
	padding: 16px 0;
}
.op-nav-link {
	display: block;
	padding: 10px 20px;
	color: #374151;
	text-decoration: none;
}
.op-nav-link:hover { background: #f0f1f4; }
.op-nav-link-active { background: #eef2ff; color: #3730a3; font-weight: 600; border-right: 3px solid #4338ca; }

.op-main { flex: 1; padding: 24px 32px; }
.op-main h1 { margin-top: 0; }

.op-page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }

.op-btn {
	display: inline-block;
	padding: 6px 14px;
	border: 1px solid #d0d3d9;
	border-radius: 6px;
	background: #ffffff;
	color: #1a1a1a;
	text-decoration: none;
	cursor: pointer;
	font-size: 14px;
}
.op-btn:hover { background: #f0f1f4; }
.op-btn-primary { background: #4338ca; border-color: #4338ca; color: #ffffff; }
.op-btn-primary:hover { background: #3730a3; }
.op-btn-danger { background: #ffffff; border-color: #dc2626; color: #dc2626; }
.op-btn-danger:hover { background: #fee2e2; }

.op-table { width: 100%; border-collapse: collapse; background: #ffffff; }
.op-table th, .op-table td { text-align: left; padding: 10px 12px; border-bottom: 1px solid #e2e4e9; }
.op-table th { font-size: 13px; text-transform: uppercase; color: #6b7280; }
.op-actions { display: flex; gap: 8px; }
.op-inline-form { margin: 0; display: inline; }
.op-empty { text-align: center; color: #6b7280; padding: 24px; }

.op-field { margin-bottom: 16px; max-width: 520px; }
.op-field label { display: block; font-weight: 600; margin-bottom: 4px; font-size: 14px; }
.op-field input[type=text], .op-field input[type=email], .op-field input[type=password], .op-field textarea {
	width: 100%;
	padding: 8px 10px;
	border: 1px solid #d0d3d9;
	border-radius: 6px;
	font-size: 14px;
}
.op-field-checkbox { display: flex; align-items: center; gap: 8px; }
.op-field-checkbox label { margin: 0; font-weight: normal; }
.op-help { color: #6b7280; font-size: 12px; margin: 4px 0 0; }

.op-error {
	background: #fef2f2;
	border: 1px solid #fecaca;
	color: #991b1b;
	padding: 10px 14px;
	border-radius: 6px;
	margin-bottom: 16px;
}
`
