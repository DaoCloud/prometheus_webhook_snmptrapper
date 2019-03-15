package snmptrapper

import (
	"time"

	types "github.com/chrusty/prometheus_webhook_snmptrapper/types"

	logrus "github.com/sirupsen/logrus"
	snmpgo "github.com/k-sone/snmpgo"
)

func sendTrap(alert types.Alert) {

	// Prepare an SNMP handler:
	snmp, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:   snmpgo.V2c,
		Address:   myConfig.SNMPTrapAddress,
		Retries:   myConfig.SNMPRetries,
		Community: myConfig.SNMPCommunity,
	})
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to create snmpgo.SNMP object")
		return
	} else {
		log.WithFields(logrus.Fields{"address": myConfig.SNMPTrapAddress, "retries": myConfig.SNMPRetries, "community": myConfig.SNMPCommunity}).Debug("Created snmpgo.SNMP object")
	}

	// Build VarBind list:
	var varBinds snmpgo.VarBinds

	// The "enterprise OID" for the trap (rising/firing or falling/recovery):
	if alert.Status == "firing" {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.FiringTrap))
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(alert.StartsAt.Format(time.RFC3339)))))
	} else {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.RecoveryTrap))
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(alert.EndsAt.Format(time.RFC3339)))))
	}

	// Insert the AlertManager variables:
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Alert, snmpgo.NewOctetString([]byte(alert.Labels["alertname"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Instance, snmpgo.NewOctetString([]byte(alert.Labels["instance"]))))
	if ns, ok := alert.Labels["namespace"]; ok {
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Tenant, snmpgo.NewOctetString([]byte(ns))))
	} else {
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Tenant, snmpgo.NewOctetString([]byte(""))))
	}

	if app, ok := alert.Labels["dce_app"]; ok {
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Application, snmpgo.NewOctetString([]byte(app))))
	} else {
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Application, snmpgo.NewOctetString([]byte(""))))
	}
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Severity, snmpgo.NewOctetString([]byte(alert.Labels["severity"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Description, snmpgo.NewOctetString([]byte(alert.Annotations["description"]))))


	// Create an SNMP "connection":
	if err = snmp.Open(); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to open SNMP connection")
		return
	}
	defer snmp.Close()

	// Send the trap:
	if err = snmp.V2Trap(varBinds); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to send SNMP trap")
		return
	} else {
		log.WithFields(logrus.Fields{"status": alert.Status}).Info("It's a trap!")
	}
}
