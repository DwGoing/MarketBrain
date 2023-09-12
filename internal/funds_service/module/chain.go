package module

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Chain struct {
	Config *Config `normal:""`
}
