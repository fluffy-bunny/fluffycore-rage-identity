package LinkedAccounts

// parseIdentityName parses an identity name string to extract provider and email
// Formats supported:
// - "google:user@gmail.com" -> ("Google", "user@gmail.com")
// - "microsoft:user@outlook.com" -> ("Microsoft", "user@outlook.com")
// - "github:username" -> ("GitHub", "username")
func parseIdentityName(name string) (provider string, email string) {
	parts := make([]string, 0)
	for i, part := range splitByColon(name) {
		if i == 0 {
			// First part is provider
			switch part {
			case "google":
				provider = "Google"
			case "microsoft":
				provider = "Microsoft"
			case "github":
				provider = "GitHub"
			case "facebook":
				provider = "Facebook"
			default:
				// Capitalize first letter
				if len(part) > 0 {
					if part[0] >= 'a' && part[0] <= 'z' {
						provider = string(part[0]-32) + part[1:]
					} else {
						provider = part
					}
				}
			}
		} else {
			parts = append(parts, part)
		}
	}

	// Join remaining parts as email/username
	if len(parts) > 0 {
		email = joinStrings(parts, ":")
	} else {
		email = name
	}
	return
}

func splitByColon(s string) []string {
	result := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			result = append(result, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func joinStrings(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	return result
}
