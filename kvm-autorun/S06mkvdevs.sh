#! /bin/sh
### BEGIN INIT INFO
# Provides:          mkvdevs
# Required-Start:    checkroot
# Required-Stop:
# Default-Start:     S
# Default-Stop:
# Short-Description: Update the /dev/vd[abcd] files
# Description:       Create the /dev/vd? device files with the 
#                    correct major number.
### END INIT INFO

PATH=/sbin:/bin:/usr/bin
. /lib/init/vars.sh
. /lib/init/tmpfs.sh

. /lib/lsb/init-functions
. /lib/init/mount-functions.sh

do_start () {

	major=$(grep virtblk /proc/devices | awk '{print $1}')

	rm -f /dev/vd?

	mknod /dev/vda b $major 0
	mknod /dev/vdb b $major 16
	mknod /dev/vdc b $major 32
	mknod /dev/vdd b $major 48
}

case "$1" in
  start|"")
	do_start
	;;
  restart|reload|force-reload)
	echo "Error: argument '$1' not supported" >&2
	exit 3
	;;
  stop)
	# No-op
	;;
  *)
	echo "Usage: mountall-mtab.sh [start|stop]" >&2
	exit 3
	;;
esac

