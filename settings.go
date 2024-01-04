package tasai

import "fmt"

type Environment string

const (
	DEV  Environment = "dev"
	PROD Environment = "prod"
)

type Instance struct {
	CPU     string
	Memory  string
	Minimum int
	Maximum int
}

type appSettings struct {
	Project   string              `json:"project"`
	Name      string              `json:"name"`
	Localport string              `json:"-"`
	Localuri  string              `json:"-"`
	Postgres  bool                `json:"postgres"`
	Instances map[string]Instance `json:"instances,omitempty"`
}

func (as *appSettings) validate() error {
	if as.Project == "" {
		return fmt.Errorf("project is required")
	}
	if as.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func defaultSettings() *appSettings {
	return &appSettings{
		Localport: ":8080",
		Localuri:  "postgres://root:password@localhost:5432/local?sslmode=disable",
		Postgres:  false,
		Instances: map[string]Instance{
			"dev": {
				CPU:     "1000m",
				Memory:  "512Mi",
				Minimum: 0,
				Maximum: 1,
			},
			"prod": {
				CPU:     "1000m",
				Memory:  "512Mi",
				Minimum: 0,
				Maximum: 1,
			},
		},
	}
}

func processAndValidateOpts(opts []appOption) (*appSettings, error) {
	s := new(appSettings)
	s.Instances = make(map[string]Instance)
	for _, opt := range opts {
		opt.apply(s)
	}
	defaults := defaultSettings()
	if s.Localport == "" {
		s.Localport = defaults.Localport
	}
	if s.Localuri == "" {
		s.Localuri = defaults.Localuri
	}
	for k, v := range s.Instances {
		if v.CPU == "" {
			v.CPU = defaults.Instances[k].CPU
		}
		if v.Memory == "" {
			v.Memory = defaults.Instances[k].Memory
		}
		if v.Maximum == 0 {
			v.Maximum = defaults.Instances[k].Maximum
		}
		s.Instances[k] = v
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return s, nil
}

type withProject string

func (w withProject) apply(a *appSettings) {
	a.Project = string(w)
}

func WithProject(project string) appOption {
	return withProject(project)
}

type withName string

func (w withName) apply(a *appSettings) {
	a.Name = string(w)
}

func WithName(name string) appOption {
	return withName(name)
}

type withLocalPort string

func (w withLocalPort) apply(a *appSettings) {
	a.Localport = string(w)
}

func WithLocalPort(port int) appOption {
	return withLocalPort(":" + fmt.Sprint(port))
}

type withPostgres []string

func (w withPostgres) apply(a *appSettings) {
	a.Postgres = true
	if len(w) > 0 {
		a.Localuri = string(w[0])
	}
}

func WithPostgres(localUri ...string) appOption {
	return withPostgres(localUri)
}

type withDevInstance Instance
type withProdInstance Instance

func (w withDevInstance) apply(a *appSettings) {
	a.Instances["dev"] = Instance(w)
}

func WithProdInstance(instance Instance) appOption {
	return withProdInstance(instance)
}

func (w withProdInstance) apply(a *appSettings) {
	a.Instances["prod"] = Instance(w)
}

func WithDevInstance(instance Instance) appOption {
	return withDevInstance(instance)
}

type appOption interface {
	apply(*appSettings)
}
