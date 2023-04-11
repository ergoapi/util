package types

type Hostinfo struct {
	OS             string `json:",omitempty"`
	OSVersion      string `json:",omitempty"`
	Container      bool   `json:",omitempty"`
	Distro         string `json:",omitempty"` // "debian", "ubuntu", "nixos", ...
	DistroVersion  string `json:",omitempty"` // "20.04", ...
	DistroCodeName string `json:",omitempty"` // "jammy", "bullseye", ...
	Machine        string `json:",omitempty"` // the current host's machine type (uname -m)
	GoArch         string `json:",omitempty"` // GOARCH value (of the built binary)
	GoArchVar      string `json:",omitempty"` // GOARM, GOAMD64, etc (of the built binary)
	GoVersion      string `json:",omitempty"` // Go version binary was built with
	Hostname       string `json:",omitempty"` // name of the host the client runs on
}
