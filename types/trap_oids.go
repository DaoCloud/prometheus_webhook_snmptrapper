package types

import (
	snmpgo "github.com/k-sone/snmpgo"
)

type TrapOIDs struct {
	FiringTrap   *snmpgo.Oid
	RecoveryTrap *snmpgo.Oid
	Instance     *snmpgo.Oid
	Severity     *snmpgo.Oid
	Description  *snmpgo.Oid
	TimeStamp    *snmpgo.Oid
}
