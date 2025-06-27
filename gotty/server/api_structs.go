package server

type ApiResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type KubeConfigRequest struct {
	Name       string `json:"name"`
	KubeConfig string `json:"kubeConfig"`
}

type KubeTokenRequest struct {
	Name      string `json:"name"`
	ApiServer string `json:"apiServer"`
	Token     string `json:"token"`
}

type FileBrowserRequest struct {
	Path  string `json:"path"`
	Name  string `json:"name,omitempty"`
	Token string `json:"token"`
}

type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
}

type FileBrowserResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Files   []FileInfo `json:"files,omitempty"`
	Path    string     `json:"path"`
}
