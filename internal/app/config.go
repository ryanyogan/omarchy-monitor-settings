package app

import (
	"github.com/ryanyogan/omarchy-monitor-settings/internal/monitor"
)

type Config struct {
	NoHyprlandCheck bool
	DebugMode       bool
	ForceLiveMode   bool
	IsTestMode      bool
}

type Services struct {
	Config          *Config
	MonitorDetector monitor.DetectorInterface
	ScalingManager  monitor.ScalingManagerInterface
	ConfigManager   monitor.ConfigManagerInterface
}

type MonitorDetectorInterface interface {
	DetectMonitors() ([]monitor.Monitor, error)
}

type ScalingManagerInterface interface {
	GetIntelligentScalingOptions(monitor monitor.Monitor) []monitor.ScalingOption
}

type ConfigManagerInterface interface {
	ApplyMonitorScale(monitor monitor.Monitor, scale float64) error
	ApplyGTKScale(scale int) error
	ApplyFontDPI(dpi int) error
	ApplyCompleteScalingOption(monitor monitor.Monitor, option monitor.ScalingOption) error
}

func NewServices(config *Config) *Services {
	return &Services{
		Config:          config,
		MonitorDetector: monitor.NewDetector(),
		ScalingManager:  monitor.NewScalingManager(),
		ConfigManager:   monitor.NewConfigManager(config.IsTestMode),
	}
}
