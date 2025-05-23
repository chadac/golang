// golang run mksyscall_zos_s390x.golang -o_sysnum zsysnum_zos_s390x.golang -o_syscall zsyscall_zos_s390x.golang -i_syscall syscall_zos_s390x.golang -o_asm zsymaddr_zos_s390x.s
// Code generated by the command above; see README.md. DO NOT EDIT.

//golang:build zos && s390x
#include "textflag.h"

//  provide the address of function variable to be fixed up.

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FlistxattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Flistxattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FremovexattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Fremovexattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FgetxattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Fgetxattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FsetxattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Fsetxattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_accept4Addr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·accept4(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_RemovexattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Removexattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_Dup3Addr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Dup3(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_DirfdAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Dirfd(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_EpollCreateAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·EpollCreate(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_EpollCreate1Addr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·EpollCreate1(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_EpollCtlAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·EpollCtl(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_EpollPwaitAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·EpollPwait(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_EpollWaitAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·EpollWait(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_EventfdAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Eventfd(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FaccessatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Faccessat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FchmodatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Fchmodat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FchownatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Fchownat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FdatasyncAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Fdatasync(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_fstatatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·fstatat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_LgetxattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Lgetxattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_LsetxattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Lsetxattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FstatfsAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Fstatfs(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FutimesAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Futimes(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_FutimesatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Futimesat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_GetrandomAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Getrandom(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_InotifyInitAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·InotifyInit(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_InotifyInit1Addr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·InotifyInit1(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_InotifyAddWatchAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·InotifyAddWatch(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_InotifyRmWatchAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·InotifyRmWatch(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_ListxattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Listxattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_LlistxattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Llistxattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_LremovexattrAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Lremovexattr(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_LutimesAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Lutimes(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_StatfsAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Statfs(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_SyncfsAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Syncfs(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_UnshareAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Unshare(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_LinkatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Linkat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_MkdiratAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Mkdirat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_MknodatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Mknodat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_PivotRootAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·PivotRoot(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_PrctlAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Prctl(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_PrlimitAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Prlimit(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_RenameatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Renameat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_Renameat2Addr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Renameat2(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_SethostnameAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Sethostname(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_SetnsAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Setns(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_SymlinkatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Symlinkat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_UnlinkatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·Unlinkat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_openatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·openat(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_openat2Addr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·openat2(SB), R8
	MOVD R8, ret+0(FP)
	RET

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

TEXT ·get_utimensatAddr(SB), NOSPLIT|NOFRAME, $0-8
	MOVD $·utimensat(SB), R8
	MOVD R8, ret+0(FP)
	RET
