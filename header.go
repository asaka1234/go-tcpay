package go_cheezeepay

func getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"charset":      "utf-8",
		//文档要求必须要有
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
	}
}
