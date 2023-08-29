package form

type CreateFolder struct {
	Title string `json:"title"`
	Root  string `json:"root"`
}

type UpdateFolder struct {
	Uid  string `json:"uid"`
	Root  string `json:"root"`
}
