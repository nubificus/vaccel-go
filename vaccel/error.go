// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#include <errno.h>
*/
import "C"

const (
	OK           int = 0              // All Good :D
	EINVAL       int = C.EINVAL       // Invalid argument
	ENOMEM       int = C.ENOMEM       // Out of memory
	ENOTSUP      int = C.ENOTSUP      // Operation not supported
	EINPROGRESS  int = C.EINPROGRESS  // Operation now in progress
	EBUSY        int = C.EBUSY        // Device or resource busy
	EEXIST       int = C.EEXIST       // File exists
	ENOENT       int = C.ENOENT       // No such file or directory
	ELIBBAD      int = C.ELIBBAD      // Corrupted shared library
	ENODEV       int = C.ENODEV       // No such device
	EIO          int = C.EIO          // I/O error
	ESESS        int = C.ECONNRESET   // Connection reset by peer
	EBACKEND     int = C.EPROTO       // Protocol error
	ENOEXEC      int = C.ENOEXEC      // Exec format error
	ENAMETOOLONG int = C.ENAMETOOLONG // File name too long
	EUSERS       int = C.EUSERS       // Too many users
	EPERM        int = C.EPERM        // Operation not permitted
	ELOOP        int = C.ELOOP        // Too many symbolic links
	EMLINK       int = C.EMLINK       // Too many links
	ENOSPC       int = C.ENOSPC       // No space left on device
	ENOTDIR      int = C.ENOTDIR      // Not a directory
	EROFS        int = C.EROFS        // Read-only file system
	EACCES       int = C.EACCES       // Permission denied
	EBADF        int = C.EBADF        // Bad file number
	EREMOTEIO    int = C.EREMOTEIO    // Remote I/O error
	EFAULT       int = C.EFAULT       // Bad address
)
