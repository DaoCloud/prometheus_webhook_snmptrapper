#!/bin/sh
/prometheus_webhook_snmptrapper -snmpcommunity=$SNMP_COMMUNITY -snmpretries=$SNMP_RETRIES -snmptrapaddress=$SNMP_TRAP_ADDRESS