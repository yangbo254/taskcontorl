package node

import (
	"encoding/json"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemInfo struct {
	CpuInfo        []cpu.InfoStat            `json:"cpuinfo"`
	CpuStat        []cpu.TimesStat           `json:"cpustat"`
	MemInfo        *mem.VirtualMemoryStat    `json:"meminfo"`
	HostInfo       *host.InfoStat            `json:"hostInfo"`
	NetConnections []net.ConnectionStat      `json:"netconnections"`
	NetIOCounters  []net.IOCountersStat      `json:"netiocounters"`
	DockerStat     []docker.CgroupDockerStat `json:"dockerstat"`
}

func (sys *SystemInfo) GetInfo() {
	sys.CpuInfo, _ = cpu.Info()
	sys.CpuStat, _ = cpu.Times(true)
	sys.MemInfo, _ = mem.VirtualMemory()
	sys.HostInfo, _ = host.Info()
	sys.NetConnections, _ = net.Connections("all")
	sys.NetIOCounters, _ = net.IOCounters(false)
	sys.DockerStat, _ = docker.GetDockerStat()

}

func (sys *SystemInfo) String() string {
	byteData, err := json.Marshal(sys)
	if err != nil {
		return ""
	}
	return string(byteData)
}
