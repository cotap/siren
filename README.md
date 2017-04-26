# Siren

Siren is a simple CLI tool for performing Consul health checks. Provided a check name along with a WARN and FAIL threshold, Siren will output the following status codes:

| Status   | Exit Code |
|:---------|:----------|
| __Ok__   | `0`       |
| __Warn__ | `1`       |
| __Fail__ | `2`       |

## Health Checks CLI

### Memory

```bash
$ siren mem 10 20
Memory usage: 78.08%

FAIL: memory usage exceeds threshold (78.08% >= 20%)
exit status 2
```

### Swap

```bash
$ siren swap 10 20
Swap usage: 88.62%

FAIL: swap usage exceeds threshold (88.62% >= 20%)
exit status 2
```

### Disk

```bash
$ siren disk 15 20
Status     Filesystem      Size Used Avail Use% Mounted on
FAIL       /dev/disk1      465G 339G  126G  73% /
OK         devfs           207K 207K    0    0% /dev
exit status 2
```

### Load

Load is based upon a normalized, 5-minute load average: The 5min load average is divided by the number of CPUs.

```bash
$ siren load 20 40
CPUs: 8
Load Averages: 2.208 1.971 1.834
Normalized Load: 27.59% 24.63% 22.93%

WARN: 5min normalized load exceeds threshold (24.63% >= 20%)
exit status 1
```

### NTP Drift

NTP takes WARN and FAIL arguments measured in milliseconds of drift.

```bash
$ siren ntp 40 100
NTP drift: 72.06ms

WARN: NTP drift exceeds threshold (72.06ms >= 40ms)
exit status 1
```

## Credits

Many thanks to the folks over at [Cloudfoundry](https://github.com/cloudfoundry) for open-sourcing their [Go Sigar port](https://github.com/cloudfoundry/gosigar) which this library relies upon!
