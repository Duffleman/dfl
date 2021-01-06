# dflmon

## Installation

- `MONITOR_CACHET_URL` - URL for the cachet status page
- `MONITOR_CACHET_KEY` - API token for the cachet status page
- `MONITOR_JOBS` - raw JSON for the jobs to handle
- `MONITOR_JOBS_FILE` - filepath to the JSON for the jobs to handle

### Config

Here is an example config. "Interval" is always in seconds. You can save the config as minified JSON as insert it as raw JSON via an environment variable called `MONITOR_JOBS`, or save it to a file and pass the path to find that file as an env var called `MONITOR_JOBS_FILE`.

```json
[
	{
		"name": "radarr",
		"component_name": "Radarr",
		"type": "https",
		"host": "radarr.int.dfl.mn",
		"interval": 5
	},
	{
		"name": "vpn_tunnel",
		"component_name": "VPN Bridge",
		"type": "icmp",
		"host": "192.168.254.254",
		"interval": 5
	}
]
```

So far only these two forms are supported:

- icmp
- https
- https-novalidate
