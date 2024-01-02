package pglookout

type ReplicationSet struct{}

type ClusterMonitor struct {
	running bool
}

// TODO clean up inputs and pass in perhaps an option struct
func NewClusterMonitor(config, cluster_state, observer_state, create_alert_file, cluster_monitor_check_queue, failover_decision_queue any, is_replication_lag_over_warning_limit bool, stats any) *ClusterMonitor {
	return nil
}

func (c *ClusterMonitor) main_monitoring_loop(requested_check bool) {}

func (c *ClusterMonitor) run() {
	c.main_monitoring_loop(false)
	for c.running {
		requested_check := false
		c.main_monitoring_loop(requested_check)
	}
}
