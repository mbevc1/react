# ReAcT
Remote Execution Automation Control Tool (Remote Execution ACTions).

## Usage
This simple tool can be used for remote execution and automation of tasks.

Example will execute command on all devices in `group1`:

```bash
# react -g group1 -c "<some_cmd>
```

Configuration of devices is read from `mts.yaml` or set by `-f` parameter.
Example configuration structure is:

```yaml
---
mts:
- host: '127.0.0.1'
  name: 'name1'
  user: 'user1'
  pass: 'pass1'
  group: 'group1'
- host: '127.0.0.2'
  name: 'name2'
  user: 'user2'
  pass: 'pass2'
  group: 'group2'
```

## Authors
* *Marko Bevc* - [@mbevc1](https://github.com/mbevc1)
