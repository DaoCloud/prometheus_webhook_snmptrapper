package snmptrapper

import (
	"github.com/chrusty/prometheus_webhook_snmptrapper/types"
	"github.com/k-sone/snmpgo"
	"github.com/sirupsen/logrus"
	"strconv"
)


func sendTraps(alerts types.AlertGroup) {
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
	//if alerts.Status == "firing" {
	//	startAt := alert.StartsAt
	//	startTime, err := startAt.MarshalText()
	//	if err != nil {
	//		log.WithFields(logrus.Fields{"error": err}).Error("Failed to marshal time to text")
	//		return
	//	}
	//	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.FiringTrap))
	//	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(startTime))))
	//} else {
	//	endAt := alert.EndsAt
	//	endTime, err := endAt.MarshalText()
	//	if err != nil {
	//		log.WithFields(logrus.Fields{"error": err}).Error("Failed to marshal time to text")
	//		return
	//	}
	//	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.RecoveryTrap))
	//	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(endTime))))
	//}

	// Insert the AlertManager variables:


	if alerts.Status == "firing" {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.FiringTrap))
	} else {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.RecoveryTrap))
	}


	var s string
	for i, alert := range alerts.Alerts {
		startAt := alert.StartsAt
		startTime := startAt.Format("2006-01-02 15:04:05")
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Error("Failed to marshal time to text")
			return
		}
		s += "No." + strconv.Itoa(i+1) + " " + alert.Annotations["description"] + ",active at " + string(startTime) + ";"
	}

	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Description, snmpgo.NewOctetString([]byte(s))))
	//varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Alert, snmpgo.NewOctetString([]byte(alert.Labels["alertname"]))))
	//varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Instance, snmpgo.NewOctetString([]byte(alert.Labels["instance"]))))
	//varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Severity, snmpgo.NewOctetString([]byte(alert.Labels["severity"]))))
	//varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Description, snmpgo.NewOctetString([]byte(alert.Annotations["description"]))))


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
		log.WithFields(logrus.Fields{"status": alerts.Status}).Info("It's a trap!")
	}
}

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
		startAt := alert.StartsAt
		startTime, err := startAt.MarshalText()
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Error("Failed to marshal time to text")
			return
		}
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.FiringTrap))
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(startTime))))
	} else {
		endAt := alert.EndsAt
		endTime, err := endAt.MarshalText()
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Error("Failed to marshal time to text")
			return
		}
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.RecoveryTrap))
		varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(endTime))))
	}

	// Insert the AlertManager variables:

	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Alert, snmpgo.NewOctetString([]byte(alert.Labels["alertname"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Instance, snmpgo.NewOctetString([]byte(alert.Labels["instance"]))))
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
