package response

type Route struct {
	Path      string  `json:"path"`
	Name      string  `json:"name"`
	Component string  `json:"component"`
	Redirect  string  `json:"redirect"`
	Meta      Meta    `json:"meta"`
	Children  []Route `json:"children"`
}

type Meta struct {
	Title              string `json:"title"`
	Icon               string `json:"icon"`
	Affix              bool   `json:"affix"`
	FrameSrc           string `json:"frameSrc"`
	IgnoreKeepAlive    bool   `json:"ignoreKeepAlive"`
	CurrentActiveMenu  string `json:"currentActiveMenu"`
	HideMenu           bool   `json:"hideMenu"`
	ShowMenu           bool   `json:"showMenu"`
	HideChildrenInMenu bool   `json:"hideChildrenInMenu"`
	HideBreadcrumb     bool   `json:"hideBreadcrumb"`
}
