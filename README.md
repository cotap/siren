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
$ siren disk 70 90
Filesystem      Size Used Avail Use% Mounted on      Status
/dev/disk1      465G 339G  126G  73% /               WARN
devfs           207K 207K    0    0% /dev            OK
exit status 1
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
