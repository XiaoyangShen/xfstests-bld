#!/bin/bash
. /root/test-config
umount $VDB
umount $VDD
/sbin/e2fsck -fy $VDB
if test $? -ge 8 ; then
	mke2fs -t ext4 $VDB
fi
dmesg -n 5
cd /root/xfstests

if test "$FSTESTCFG" = all
then
	FSTESTCFG="4k ext3 nojournal 1k ext3conv metacsum dioread_nolock data_journal bigalloc"
fi

export SCRATCH_DEV=$VDC
export SCRATCH_MNT=/vdc
for i in $FSTESTCFG
do
	if test -e "/root/conf/$i"; then
		. /root/conf/$i
	else
		echo "Unknown configuration $i!"
		continue
	fi
	if test "$TEST_DEV" = "$VDD" ; then
		if test "$FS" = "ext4" ; then
		    mke2fs -q -t ext4 $MKFS_OPTIONS $TEST_DEV
		elif test "$FS" = "xfs" ; then
		    /sbin/mkfs.xfs -f $MKFS_OPTIONS $TEST_DEV
		else
		    /sbin/mkfs.$FS $TEST_DEV
		fi
	fi
	echo -n "BEGIN TEST: $TESTNAME " ; date
	echo Device: $TEST_DEV
	echo mk2fs options: $MKFS_OPTIONS
	echo mount options: $EXT_MOUNT_OPTIONS
	sh ./check -ext4 $FSTESTSET
	echo -n "END TEST: $TESTNAME " ; date
	umount $TEST_DEV >& /dev/null
	if test "$FS" = "ext4" ; then
		/sbin/e2fsck -fy $TEST_DEV
	elif test "$FS" = "xfs" ; then
		xfs_check -f $TEST_DEV
	else
		/sbin/fsck.$FS $TEST_DEV
	fi
done
