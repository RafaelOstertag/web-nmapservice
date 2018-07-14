#!/bin/sh
#
# Use this mock when running nmap_test.go. It can be specified
# by setting the environment variable NMAP_CMD.

EXPECTED_ARGS="-oX - -p 22 gizmo.kruemel.home"

if [ $# -ne 5 ]
then
    echo "Wrong arguments: $@" >&2
    echo "Expected: ${EXPECTED_ARGS}" >&2
    exit 1
fi

for arg in ${EXPECTED_ARGS}
do
    if test "$1" != "${arg}"
    then
        echo "Expected '$1' to be '${arg}'" >&2
        exit 1
    fi
    shift
done

cat <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE nmaprun>
<?xml-stylesheet href="file:///usr/bin/../share/nmap/nmap.xsl" type="text/xsl"?>
<!-- Nmap 7.60 scan initiated Thu Jul 12 20:55:19 2018 as: nmap -oX - -p 22 gizmo.kruemel.home -->
<nmaprun scanner="nmap" args="nmap -oX - -p 22 gizmo.kruemel.home" start="1531421719" startstr="Thu Jul 12 20:55:19 2018" version="7.60" xmloutputversion="1.04">
<scaninfo type="connect" protocol="tcp" numservices="1" services="22"/>
<verbose level="0"/>
<debugging level="0"/>
<host starttime="1531421719" endtime="1531421719"><status state="up" reason="syn-ack" reason_ttl="0"/>
<address addr="192.168.100.1" addrtype="ipv4"/>
<hostnames>
<hostname name="gizmo.kruemel.home" type="user"/>
</hostnames>
<ports><port protocol="tcp" portid="22"><state state="open" reason="syn-ack" reason_ttl="0"/><service name="ssh" method="table" conf="3"/></port>
</ports>
<times srtt="4247" rttvar="4323" to="100000"/>
</host>
<runstats><finished time="1531421719" timestr="Thu Jul 12 20:55:19 2018" elapsed="0.06" summary="Nmap done at Thu Jul 12 20:55:19 2018; 1 IP address (1 host up) scanned in 0.06 seconds" exit="success"/><hosts up="1" down="0" total="1"/>
</runstats>
</nmaprun>
EOF

exit 0