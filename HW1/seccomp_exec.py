import seccomp, errno, os, sys
mydocker = seccomp.SyscallFilter(seccomp.ALLOW)
mydocker.add_rule(seccomp.ERRNO(errno.EPERM), "mkdir")
mydocker.add_rule(seccomp.ERRNO(errno.EPERM), "mkdirat")
mydocker.load()
os.execvp(sys.argv[1], sys.argv[1:])
