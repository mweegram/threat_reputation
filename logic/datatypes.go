package logic

type Threat struct {
	ID        int
	Filename  string
	Sha256    string
	Comments  []string
	Submitted string
}

type DisplayThreat struct {
	ID        int
	Filename  string
	Sha256    string
	Comments  []Comment
	Submitted string
}

type Homepage_Content struct {
	Threats []Homepage_Threat
	Stats   []Stats
}

type Homepage_Threat struct {
	ID       int
	Filename string
}

type Comment struct {
	ID   int
	Text string
	Date string
}

type Submission struct {
	Filename string
	Filehash string
}

type Notification struct {
	Status  string
	Message string
	Colour  string
}

type Stats struct {
	Title string
	Value int
}

type Search_Result struct {
	ID       int
	Filename string
	Sha256   string
}
