package executors

type PackageInfo struct {
	Name    string
	Version string
	Arch    string
}

type Executor interface {
	GetPackages() ([]*PackageInfo, error)
	Update(string) error
	BulkUpdate([]string) error
}
