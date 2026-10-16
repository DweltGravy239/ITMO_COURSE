#!/bin/bash
set -euo pipefail
memory_max="${memory_max:-100M}"
pids_max="${pids_max:-100}"
if [[ "$memory_max" =~ ^[0-9]+[KMG]?$ ]]; then
	if [[ "$pids_max" =~ ^[0-9]+$ ]]; then
		mkdir -p /sys/fs/cgroup/mydocker
		echo $memory_max > /sys/fs/cgroup/mydocker/memory.max
		echo 0 > /sys/fs/cgroup/mydocker/memory.swap.max
		echo '50000 100000' > /sys/fs/cgroup/mydocker/cpu.max
		echo $pids_max > /sys/fs/cgroup/mydocker/pids.max
		echo $$ > /sys/fs/cgroup/mydocker/cgroup.procs
		cat /sys/fs/cgroup/mydocker/memory.max
		cat /sys/fs/cgroup/mydocker/pids.max
		cat /proc/self/cgroup
		unshare --pid --fork --mount-proc --uts --net --ipc /bin/bash -c 'hostname myapp; ip link set lo up; exec capsh --drop=cap_sys_time,cap_sys_admin -- -c "exec python3 /ITMO_COURSE/HW1/seccomp_exec.py /ITMO_COURSE/HW1/api/api"'
	else
		echo "Вы ввели не число"
	fi
else
	echo "Вы ввели не число"
fi
