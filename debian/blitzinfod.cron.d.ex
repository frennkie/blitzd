#
# Regular cron jobs for the blitzinfod package
#
0 4	* * *	root	[ -x /usr/bin/blitzinfod_maintenance ] && /usr/bin/blitzinfod_maintenance
