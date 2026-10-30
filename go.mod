module golang.zabbix.com/agent2/zbxedx

go 1.25.12

require (
	golang.zabbix.com/agent2 v0.0.0
	golang.zabbix.com/sdk v1.2.2-0.20261027065228-4735c97f7e8b
)

require golang.org/x/sys v0.47.0 // indirect

replace golang.zabbix.com/agent2 => ./zabbix/src/go
