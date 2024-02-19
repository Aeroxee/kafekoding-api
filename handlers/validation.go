package handlers

func isAllowedExtension(ext string) bool {
	allowedExtension := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	return allowedExtension[ext]
}
