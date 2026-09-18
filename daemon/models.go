package main

type RAMInfo struct {
	TotalKB uint64 `json:"total_kb"`
	FreeKB  uint64 `json:"free_kb"`
	UsedKB  uint64 `json:"used_kb"`
}

type ProcessInfo struct {
	PID            int    `json:"pid"`
	Name           string `json:"name"`
	Cmd            string `json:"cmd"`
	VSZKB          uint64 `json:"vsz_kb"`
	RSSKB          uint64 `json:"rss_kb"`
	MemPercentX100 uint64 `json:"mem_percent_x100"`
	CPUPercentX100 uint64 `json:"cpu_percent_x100"`
}

type KernelInfo struct {
	RAM       RAMInfo       `json:"ram"`
	Processes []ProcessInfo `json:"processes"`
}

type ContainerInfo struct {
	ID       string
	Name     string
	PID      int
	Profile  string
	Resource string
	Process  *ProcessInfo
}
