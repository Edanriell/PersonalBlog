package utils

func MarkdownToHTML(markdown string) string {
	return string(blackfriday.Run([]byte(markdown)))
}
