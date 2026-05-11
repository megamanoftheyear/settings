package tree

type Meta map[string]any

type Link struct {
	Path    string
	Reason  string
	Context string
	UserID  string
}
